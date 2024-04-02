package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/app/model"
	"main/app/shared/database"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func validateTask(task model.Task) error {
	if len(task.Date) == 0 {
		task.Date = time.Now().Format(model.DatePattern)
	}

	date, err := time.Parse(model.DatePattern, task.Date)
	if err != nil {
		return err
	}

	if date.Before(time.Now()) {
		task.Date = time.Now().Format(model.DatePattern)
	}

	if len(task.Title) == 0 {
		return errors.New("title can not be empty")
	}

	if len(task.Repeat) > 0 {
		dayRepeatRule, _ := regexp.MatchString(`d \d+`, task.Repeat)
		yearRepeatRule, _ := regexp.MatchString(`y`, task.Repeat)

		if !dayRepeatRule && !yearRepeatRule {
			return errors.New("bad repeat format")
		}
	}

	return nil
}

func TaskAddPOST(w http.ResponseWriter, r *http.Request) {
	var taskData model.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		http.Error(w, fmt.Errorf("body getting error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		http.Error(w, fmt.Errorf("JSON encoding error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	err := validateTask(taskData)
	if err != nil {
		taskIdData, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%w", err).Error()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(taskIdData)
		if err != nil {
			http.Error(w, fmt.Errorf("%w", err).Error(), http.StatusBadRequest)
		}
	}

	taskId, err := database.InsertTask(taskData)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
		return
	}

	taskIdData, err := json.Marshal(model.TaskIdResponse{Id: taskId})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(taskIdData)
	if err != nil {
		http.Error(w, fmt.Errorf("writing task id error: %w", err).Error(), http.StatusBadRequest)
	}
}

func TasksReadGET(w http.ResponseWriter, _ *http.Request) {
	tasks, err := database.ReadTasks()
	if err != nil {
		http.Error(w, fmt.Errorf("writing task id error: %w", err).Error(), http.StatusBadRequest)
	}

	tasksData, err := json.Marshal(model.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
}

func TaskReadGET(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := database.ReadTask(id)
	if err != nil {
		errorData, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%w", err).Error()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorData)
		return
	}

	taskData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)
}

func TaskUpdatePUT(w http.ResponseWriter, r *http.Request) {
	var taskData model.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		http.Error(w, fmt.Errorf("body getting error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		http.Error(w, fmt.Errorf("JSON encoding error: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(taskData.Id); err != nil {
		errorData, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%w", err).Error()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorData)
		return
	}

	err := validateTask(taskData)
	if err != nil {
		errorData, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%w", err).Error()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorData)
		return
	}

	affected, err := database.UpdateTask(taskData)
	if affected == 0 || err != nil {
		errorData, _ := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("%w", err).Error()})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorData)
		return
	}

	taskIdData, err := json.Marshal(taskData)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskIdData)
	if err != nil {
		http.Error(w, fmt.Errorf("writing task id error: %w", err).Error(), http.StatusBadRequest)
	}
}
