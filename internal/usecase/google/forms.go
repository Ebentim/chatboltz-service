package usecase

import (
	gforms "google.golang.org/api/forms/v1"
)

// FormsIntegration defines the interface for forms operations.
type FormsIntegration interface {
	CreateForm(*gforms.Form) (*gforms.Form, error)
	GetForm(formID string) (*gforms.Form, error)
	GetFormResponses(formID string) (*gforms.ListFormResponsesResponse, error)
	BatchUpdate(formID string, requests []*gforms.Request) (*gforms.BatchUpdateFormResponse, error)
}

// FormsUseCase handles the business logic for forms operations.
type FormsUseCase struct {
	formsIntegration FormsIntegration
}

// NewFormsUseCase creates a new FormsUseCase.
func NewFormsUseCase(fi FormsIntegration) *FormsUseCase {
	return &FormsUseCase{formsIntegration: fi}
}

// CreateForm creates a new Google Form.
func (uc *FormsUseCase) CreateForm(title, description string) (*gforms.Form, error) {
	form := &gforms.Form{
		Info: &gforms.Info{
			Title:       title,
			Description: description,
		},
	}
	return uc.formsIntegration.CreateForm(form)
}

// GetForm retrieves a Google Form.
func (uc *FormsUseCase) GetForm(formID string) (*gforms.Form, error) {
	return uc.formsIntegration.GetForm(formID)
}

// GetFormResponses retrieves all responses from a Google Form.
func (uc *FormsUseCase) GetFormResponses(formID string) (*gforms.ListFormResponsesResponse, error) {
	return uc.formsIntegration.GetFormResponses(formID)
}

// MakeQuiz turns a form into a quiz.
func (uc *FormsUseCase) MakeQuiz(formID string) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			UpdateSettings: &gforms.UpdateSettingsRequest{
				Settings: &gforms.FormSettings{
					QuizSettings: &gforms.QuizSettings{
						IsQuiz: true,
					},
				},
				UpdateMask: "quizSettings.isQuiz",
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// --- Question Types ---

// AddTextQuestion adds a short answer text question.
func (uc *FormsUseCase) AddTextQuestion(formID, title string, location int64, points int64, correctAnswer string) (*gforms.BatchUpdateFormResponse, error) {
	question := &gforms.Question{
		TextQuestion: &gforms.TextQuestion{},
	}
	if points > 0 {
		question.Grading = &gforms.Grading{
			PointValue: points,
			CorrectAnswers: &gforms.CorrectAnswers{
				Answers: []*gforms.CorrectAnswer{{Value: correctAnswer}},
			},
		}
	}

	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: question,
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddParagraphQuestion adds a long answer text question.
func (uc *FormsUseCase) AddParagraphQuestion(formID, title string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: &gforms.Question{
							TextQuestion: &gforms.TextQuestion{Paragraph: true},
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddMultipleChoiceQuestion adds a multiple choice question.
func (uc *FormsUseCase) AddMultipleChoiceQuestion(formID, title string, choices []string, location int64, points int64, correctAnswer string) (*gforms.BatchUpdateFormResponse, error) {
	var options []*gforms.Option
	for _, choice := range choices {
		options = append(options, &gforms.Option{Value: choice})
	}
	question := &gforms.Question{
		ChoiceQuestion: &gforms.ChoiceQuestion{
			Type:    "RADIO",
			Options: options,
		},
	}
	if points > 0 {
		question.Grading = &gforms.Grading{
			PointValue: points,
			CorrectAnswers: &gforms.CorrectAnswers{
				Answers: []*gforms.CorrectAnswer{{Value: correctAnswer}},
			},
		}
	}
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: question,
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddCheckboxQuestion adds a checkbox question.
func (uc *FormsUseCase) AddCheckboxQuestion(formID, title string, choices []string, location int64, points int64, correctAnswers []string) (*gforms.BatchUpdateFormResponse, error) {
	var options []*gforms.Option
	for _, choice := range choices {
		options = append(options, &gforms.Option{Value: choice})
	}
	question := &gforms.Question{
		ChoiceQuestion: &gforms.ChoiceQuestion{
			Type:    "CHECKBOX",
			Options: options,
		},
	}
	if points > 0 && len(correctAnswers) > 0 {
		var answers []*gforms.CorrectAnswer
		for _, ans := range correctAnswers {
			answers = append(answers, &gforms.CorrectAnswer{Value: ans})
		}
		question.Grading = &gforms.Grading{
			PointValue: points,
			CorrectAnswers: &gforms.CorrectAnswers{
				Answers: answers,
			},
		}
	}
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: question,
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddDropdownQuestion adds a dropdown question.
func (uc *FormsUseCase) AddDropdownQuestion(formID, title string, choices []string, location int64, points int64, correctAnswer string) (*gforms.BatchUpdateFormResponse, error) {
	var options []*gforms.Option
	for _, choice := range choices {
		options = append(options, &gforms.Option{Value: choice})
	}
	question := &gforms.Question{
		ChoiceQuestion: &gforms.ChoiceQuestion{
			Type:    "DROP_DOWN",
			Options: options,
		},
	}
	if points > 0 {
		question.Grading = &gforms.Grading{
			PointValue: points,
			CorrectAnswers: &gforms.CorrectAnswers{
				Answers: []*gforms.CorrectAnswer{{Value: correctAnswer}},
			},
		}
	}
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: question,
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddFileUploadQuestion adds a file upload question.
func (uc *FormsUseCase) AddFileUploadQuestion(formID, title, folderId string, types []string, maxFiles int64, maxSize int64, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: &gforms.Question{
							FileUploadQuestion: &gforms.FileUploadQuestion{
								FolderId:    folderId,
								Types:       types,
								MaxFiles:    maxFiles,
								MaxFileSize: maxSize,
							},
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddLinearScaleQuestion adds a linear scale question.
func (uc *FormsUseCase) AddLinearScaleQuestion(formID, title string, low, high int64, lowLabel, highLabel string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: &gforms.Question{
							ScaleQuestion: &gforms.ScaleQuestion{
								Low:       low,
								High:      high,
								LowLabel:  lowLabel,
								HighLabel: highLabel,
							},
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddMultipleChoiceGridQuestion adds a multiple choice grid question.
func (uc *FormsUseCase) AddMultipleChoiceGridQuestion(formID, title string, rows []string, columns []string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionGroupItem: &gforms.QuestionGroupItem{
						Grid: &gforms.Grid{
							Columns: &gforms.ChoiceQuestion{
								Type: "RADIO",
								Options: func() []*gforms.Option {
									var options []*gforms.Option
									for _, c := range columns {
										options = append(options, &gforms.Option{Value: c})
									}
									return options
								}(),
							},
						},
						Questions: func() []*gforms.Question {
							var questions []*gforms.Question
							for _, r := range rows {
								questions = append(questions, &gforms.Question{RowQuestion: &gforms.RowQuestion{Title: r}})
							}
							return questions
						}(),
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddCheckboxGridQuestion adds a checkbox grid question.
func (uc *FormsUseCase) AddCheckboxGridQuestion(formID, title string, rows []string, columns []string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionGroupItem: &gforms.QuestionGroupItem{
						Grid: &gforms.Grid{
							Columns: &gforms.ChoiceQuestion{
								Type: "CHECKBOX",
								Options: func() []*gforms.Option {
									var options []*gforms.Option
									for _, c := range columns {
										options = append(options, &gforms.Option{Value: c})
									}
									return options
								}(),
							},
						},
						Questions: func() []*gforms.Question {
							var questions []*gforms.Question
							for _, r := range rows {
								questions = append(questions, &gforms.Question{RowQuestion: &gforms.RowQuestion{Title: r}})
							}
							return questions
						}(),
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddDateQuestion adds a date question.
func (uc *FormsUseCase) AddDateQuestion(formID, title string, includeTime, includeYear bool, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: &gforms.Question{
							DateQuestion: &gforms.DateQuestion{
								IncludeTime: includeTime,
								IncludeYear: includeYear,
							},
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddTimeQuestion adds a time question.
func (uc *FormsUseCase) AddTimeQuestion(formID, title string, duration bool, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title: title,
					QuestionItem: &gforms.QuestionItem{
						Question: &gforms.Question{
							TimeQuestion: &gforms.TimeQuestion{
								Duration: duration,
							},
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// --- Extra Features ---

// AddSectionBreak adds a section break.
func (uc *FormsUseCase) AddSectionBreak(formID, title, description string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title:         title,
					Description:   description,
					PageBreakItem: &gforms.PageBreakItem{},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddImageItem adds an image to the form.
func (uc *FormsUseCase) AddImageItem(formID, imageUri, altText string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					ImageItem: &gforms.ImageItem{
						Image: &gforms.Image{
							SourceUri: imageUri,
							Properties: &gforms.MediaProperties{
								Alignment: "CENTER",
							},
							AltText: altText,
						},
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddVideoItem adds a video to the form.
func (uc *FormsUseCase) AddVideoItem(formID, youtubeUri, caption string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					VideoItem: &gforms.VideoItem{
						Video: &gforms.Video{
							YoutubeUri: youtubeUri,
							Properties: &gforms.MediaProperties{
								Alignment: "CENTER",
							},
						},
						Caption: caption,
					},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}

// AddTextItem adds a descriptive text block.
func (uc *FormsUseCase) AddTextItem(formID, title, description string, location int64) (*gforms.BatchUpdateFormResponse, error) {
	requests := []*gforms.Request{
		{
			CreateItem: &gforms.CreateItemRequest{
				Item: &gforms.Item{
					Title:       title,
					Description: description,
					TextItem:    &gforms.TextItem{},
				},
				Location: &gforms.Location{Index: location},
			},
		},
	}
	return uc.formsIntegration.BatchUpdate(formID, requests)
}
