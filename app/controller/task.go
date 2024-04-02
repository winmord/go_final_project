package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/app/model"
	"main/app/shared/database"
	"net/http"
	"regexp"
	"time"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
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

	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(model.DatePattern)
	}

	date, err := time.Parse(model.DatePattern, taskData.Date)
	if err != nil {
		errorData, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("bad date format: %w", err).Error()})
		if err != nil {
			http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(errorData)
		return
	}

	if date.Before(time.Now()) {
		taskData.Date = time.Now().Format(model.DatePattern)
	}

	if len(taskData.Title) == 0 {
		errorResponse, err := json.Marshal(model.ErrorResponse{Error: "title can not be empty"})
		if err != nil {
			http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
			return
		}

		_, err = w.Write(errorResponse)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
			return
		}
		return
	}

	taskId, err := database.InsertTask(taskData)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if len(taskData.Repeat) > 0 {
		dayRepeatRule, err := regexp.MatchString(`d \d+`, taskData.Repeat)
		yearRepeatRule, err := regexp.MatchString(`y`, taskData.Repeat)

		if !dayRepeatRule && !yearRepeatRule {
			errorResponse, err := json.Marshal(model.ErrorResponse{Error: fmt.Errorf("bad repeat rule: %w", err).Error()})
			if err != nil {
				http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
				return
			}

			_, err = w.Write(errorResponse)
			if err != nil {
				http.Error(w, fmt.Errorf("failed to create task: %w", err).Error(), http.StatusBadRequest)
				return
			}
			return
		}
	}

	taskIdData, err := json.Marshal(model.TaskIdResponse{Id: taskId})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(taskIdData)

	if err != nil {
		http.Error(w, fmt.Errorf("writing task id error: %w", err).Error(), http.StatusBadRequest)
	}
}
