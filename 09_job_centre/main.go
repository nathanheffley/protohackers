package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sort"
)

type RequestType struct {
	Request string
}

type PutRequest struct {
	Queue string
	Job   json.RawMessage
	Pri   int
}

type GetRequest struct {
	Queues []string
}

type DeleteRequest struct {
	Id int
}

type AbortRequest struct {
	Id int
}

type Job struct {
	id         int
	json       json.RawMessage
	pri        int
	heldByConn *net.Conn
	deleted    bool
}

type Queue struct {
	jobs []*Job
}

func (q *Queue) put(job *Job) {
	q.jobs = append(q.jobs, job)
	sort.Slice(q.jobs, func(i, j int) bool {
		return q.jobs[i].pri < q.jobs[j].pri
	})
}

func (q *Queue) get() *Job {
	if len(q.jobs) == 0 {
		return nil
	}

	for _, job := range q.jobs {
		if job.heldByConn == nil && !job.deleted {
			return job
		}
	}

	return nil
}

var queues = make(map[string]*Queue)

var jobs = make(map[int]*Job)

var uniqueIdChan = make(chan int)

func getUniqueId() int {
	return <-uniqueIdChan
}

func main() {
	listener, _ := net.Listen("tcp", "0.0.0.0:8080")
	defer listener.Close()

	go func() {
		nextId := 1
		for {
			uniqueIdChan <- nextId
			nextId++
		}
	}()

	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	for {
		reqBytes := make([]byte, 1000)
		_, err := conn.Read(reqBytes)
		if err != nil {
			if err.Error() == "EOF" {
				for _, job := range jobs {
					if job.heldByConn == &conn {
						job.heldByConn = nil
					}
				}
				break
			}

			fmt.Println("Could not read request: ", err)
			conn.Write([]byte("{\"status\":\"error\",\"error\":\"Could not read request.\"}"))
			continue
		}

		reqBytes = bytes.Trim(reqBytes, "\x00")
		reqBytes = bytes.TrimSpace(reqBytes)
		fmt.Println(string(reqBytes))
		var req RequestType
		err = json.Unmarshal(reqBytes, &req)
		if err != nil {
			fmt.Println("Invalid request format: ", err)
			conn.Write([]byte("{\"status\":\"error\",\"error\":\"Invalid request format.\"}"))
			continue
		}

		if req.Request == "put" {
			var putReq PutRequest
			err = json.Unmarshal(reqBytes, &putReq)
			if err != nil {
				fmt.Println("Could not decode full put request: ", err)
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Could not decode full put request.\"}"))
				continue
			}

			if json.Valid([]byte(putReq.Job)) == false {
				fmt.Println("Invalid JSON for job: ", err)
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Invalid JSON for job.\"}"))
				continue
			}

			queue, ok := queues[putReq.Queue]
			if !ok {
				queue = &Queue{
					jobs: []*Job{},
				}
				queues[putReq.Queue] = queue
			}

			job := Job{
				id:         getUniqueId(),
				json:       putReq.Job,
				pri:        putReq.Pri,
				heldByConn: nil,
				deleted:    false,
			}

			queue.put(&job)
			jobs[job.id] = &job

			conn.Write([]byte(fmt.Sprintf("{\"status\":\"ok\",\"id\":%d}", job.id)))

			continue
		}

		if req.Request == "get" {
			var getReq GetRequest
			err = json.Unmarshal(reqBytes, &getReq)
			if err != nil {
				fmt.Println("Could not decode full get request: ", err)
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Could not decode full get request.\"}"))
				continue
			}

			var highestPriJob *Job
			var highestPriJobQueue string
			for _, queueName := range getReq.Queues {
				queue, ok := queues[queueName]
				if !ok {
					continue
				}

				job := queue.get()
				if job == nil {
					continue
				}

				if highestPriJob == nil || job.pri > highestPriJob.pri {
					highestPriJob = job
					highestPriJobQueue = queueName
				}
			}

			if highestPriJob == nil {
				conn.Write([]byte("{\"status\":\"no-job\"}"))
				continue
			}

			highestPriJob.heldByConn = &conn

			conn.Write([]byte(fmt.Sprintf(
				"{\"status\":\"ok\",\"id\":%d,\"job\":%s,\"pri\":%d,\"queue\":\"%s\"}",
				highestPriJob.id,
				highestPriJob.json,
				highestPriJob.pri,
				highestPriJobQueue,
			)))

			continue
		}

		if req.Request == "delete" {
			var deleteReq DeleteRequest
			err = json.Unmarshal(reqBytes, &deleteReq)
			if err != nil {
				fmt.Println("Could not decode full delete request: ", err)
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Could not decode full delete request.\"}"))
				continue
			}

			job, ok := jobs[deleteReq.Id]
			if !ok {
				conn.Write([]byte("{\"status\":\"no-job\"}"))
				continue
			}

			if job.deleted {
				conn.Write([]byte("{\"status\":\"no-job\"}"))
				continue
			}

			job.deleted = true

			conn.Write([]byte("{\"status\":\"ok\"}"))

			continue
		}

		if req.Request == "abort" {
			var abortReq AbortRequest
			err = json.Unmarshal(reqBytes, &abortReq)
			if err != nil {
				fmt.Println("Could not decode full abort request: ", err)
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Could not decode full abort request.\"}"))
				continue
			}

			job, ok := jobs[abortReq.Id]
			if !ok {
				conn.Write([]byte("{\"status\":\"no-job\"}"))
				continue
			}

			if job.deleted {
				conn.Write([]byte("{\"status\":\"no-job\"}"))
				continue
			}

			if job.heldByConn != &conn {
				conn.Write([]byte("{\"status\":\"error\",\"error\":\"Cannot abort a job claimed by another client.\"}"))
				continue
			}

			job.heldByConn = nil

			conn.Write([]byte("{\"status\":\"ok\"}"))

			continue
		}

		fmt.Println("Unrecognised request type: ", req.Request)
		conn.Write([]byte("{\"status\":\"error\",\"error\":\"Unrecognised request type.\"}"))
	}
}
