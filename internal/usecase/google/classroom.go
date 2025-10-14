package usecase

import (
	gclassroom "google.golang.org/api/classroom/v1"
)

// ClassroomIntegration defines the interface for classroom operations.
type ClassroomIntegration interface {
	// Courses
	CreateCourse(course *gclassroom.Course) (*gclassroom.Course, error)
	GetCourse(courseID string) (*gclassroom.Course, error)
	ListCourses() (*gclassroom.ListCoursesResponse, error)
	UpdateCourse(courseID string, course *gclassroom.Course) (*gclassroom.Course, error)
	DeleteCourse(courseID string) error
	// Aliases
	CreateAlias(courseID string, alias *gclassroom.CourseAlias) (*gclassroom.CourseAlias, error)
	ListAliases(courseID string) (*gclassroom.ListCourseAliasesResponse, error)
	DeleteAlias(courseID, alias string) error
	// Students
	CreateStudent(courseID string, student *gclassroom.Student) (*gclassroom.Student, error)
	GetStudent(courseID, userID string) (*gclassroom.Student, error)
	ListStudents(courseID string) (*gclassroom.ListStudentsResponse, error)
	DeleteStudent(courseID, userID string) error
	// Teachers
	CreateTeacher(courseID string, teacher *gclassroom.Teacher) (*gclassroom.Teacher, error)
	GetTeacher(courseID, userID string) (*gclassroom.Teacher, error)
	ListTeachers(courseID string) (*gclassroom.ListTeachersResponse, error)
	DeleteTeacher(courseID, userID string) error
	// Invitations
	CreateInvitation(invitation *gclassroom.Invitation) (*gclassroom.Invitation, error)
	GetInvitation(id string) (*gclassroom.Invitation, error)
	ListInvitations() (*gclassroom.ListInvitationsResponse, error)
	DeleteInvitation(id string) error
	AcceptInvitation(id string) error
	// User Profiles
	GetUserProfile(userID string) (*gclassroom.UserProfile, error)
	// Topics
	CreateTopic(courseID string, topic *gclassroom.Topic) (*gclassroom.Topic, error)
	GetTopic(courseID, topicID string) (*gclassroom.Topic, error)
	ListTopics(courseID string) (*gclassroom.ListTopicResponse, error)
	UpdateTopic(courseID, topicID string, topic *gclassroom.Topic) (*gclassroom.Topic, error)
	DeleteTopic(courseID, topicID string) error
	// CourseWork
	CreateCourseWork(courseID string, courseWork *gclassroom.CourseWork) (*gclassroom.CourseWork, error)
	GetCourseWork(courseID, courseWorkID string) (*gclassroom.CourseWork, error)
	ListCourseWork(courseID string) (*gclassroom.ListCourseWorkResponse, error)
	UpdateCourseWork(courseID, courseWorkID string, courseWork *gclassroom.CourseWork) (*gclassroom.CourseWork, error)
	DeleteCourseWork(courseID, courseWorkID string) error
	// StudentSubmissions
	GetStudentSubmission(courseID, courseWorkID, submissionID string) (*gclassroom.StudentSubmission, error)
	ListStudentSubmissions(courseID, courseWorkID string) (*gclassroom.ListStudentSubmissionsResponse, error)
	PatchStudentSubmission(courseID, courseWorkID, submissionID string, submission *gclassroom.StudentSubmission, updateMask string) (*gclassroom.StudentSubmission, error)
	TurnInStudentSubmission(courseID, courseWorkID, submissionID string) error
	ReturnStudentSubmission(courseID, courseWorkID, submissionID string) error
	ReclaimStudentSubmission(courseID, courseWorkID, submissionID string) error
	// Announcements
	CreateAnnouncement(courseID string, announcement *gclassroom.Announcement) (*gclassroom.Announcement, error)
	GetAnnouncement(courseID, announcementID string) (*gclassroom.Announcement, error)
	ListAnnouncements(courseID string) (*gclassroom.ListAnnouncementsResponse, error)
	UpdateAnnouncement(courseID, announcementID string, announcement *gclassroom.Announcement) (*gclassroom.Announcement, error)
	DeleteAnnouncement(courseID, announcementID string) error
	// CourseWorkMaterials
	CreateCourseWorkMaterial(courseID string, material *gclassroom.CourseWorkMaterial) (*gclassroom.CourseWorkMaterial, error)
	GetCourseWorkMaterial(courseID, materialID string) (*gclassroom.CourseWorkMaterial, error)
	ListCourseWorkMaterials(courseID string) (*gclassroom.ListCourseWorkMaterialResponse, error)
	UpdateCourseWorkMaterial(courseID, materialID string, material *gclassroom.CourseWorkMaterial) (*gclassroom.CourseWorkMaterial, error)
	DeleteCourseWorkMaterial(courseID, materialID string) error
}

// ClassroomUseCase handles the business logic for classroom operations.
type ClassroomUseCase struct {
	classroomIntegration ClassroomIntegration
}

// NewClassroomUseCase creates a new ClassroomUseCase.
func NewClassroomUseCase(ci ClassroomIntegration) *ClassroomUseCase {
	return &ClassroomUseCase{classroomIntegration: ci}
}

// --- Courses ---
func (uc *ClassroomUseCase) CreateCourse(course *gclassroom.Course) (*gclassroom.Course, error) {
	return uc.classroomIntegration.CreateCourse(course)
}
func (uc *ClassroomUseCase) GetCourse(courseID string) (*gclassroom.Course, error) {
	return uc.classroomIntegration.GetCourse(courseID)
}
func (uc *ClassroomUseCase) ListCourses() (*gclassroom.ListCoursesResponse, error) {
	return uc.classroomIntegration.ListCourses()
}
func (uc *ClassroomUseCase) UpdateCourse(courseID string, course *gclassroom.Course) (*gclassroom.Course, error) {
	return uc.classroomIntegration.UpdateCourse(courseID, course)
}
func (uc *ClassroomUseCase) DeleteCourse(courseID string) error {
	return uc.classroomIntegration.DeleteCourse(courseID)
}

// --- Aliases ---
func (uc *ClassroomUseCase) CreateAlias(courseID string, alias *gclassroom.CourseAlias) (*gclassroom.CourseAlias, error) {
	return uc.classroomIntegration.CreateAlias(courseID, alias)
}
func (uc *ClassroomUseCase) ListAliases(courseID string) (*gclassroom.ListCourseAliasesResponse, error) {
	return uc.classroomIntegration.ListAliases(courseID)
}
func (uc *ClassroomUseCase) DeleteAlias(courseID, alias string) error {
	return uc.classroomIntegration.DeleteAlias(courseID, alias)
}

// --- Students ---
func (uc *ClassroomUseCase) CreateStudent(courseID string, student *gclassroom.Student) (*gclassroom.Student, error) {
	return uc.classroomIntegration.CreateStudent(courseID, student)
}
func (uc *ClassroomUseCase) GetStudent(courseID, userID string) (*gclassroom.Student, error) {
	return uc.classroomIntegration.GetStudent(courseID, userID)
}
func (uc *ClassroomUseCase) ListStudents(courseID string) (*gclassroom.ListStudentsResponse, error) {
	return uc.classroomIntegration.ListStudents(courseID)
}
func (uc *ClassroomUseCase) DeleteStudent(courseID, userID string) error {
	return uc.classroomIntegration.DeleteStudent(courseID, userID)
}

// --- Teachers ---
func (uc *ClassroomUseCase) CreateTeacher(courseID string, teacher *gclassroom.Teacher) (*gclassroom.Teacher, error) {
	return uc.classroomIntegration.CreateTeacher(courseID, teacher)
}
func (uc *ClassroomUseCase) GetTeacher(courseID, userID string) (*gclassroom.Teacher, error) {
	return uc.classroomIntegration.GetTeacher(courseID, userID)
}
func (uc *ClassroomUseCase) ListTeachers(courseID string) (*gclassroom.ListTeachersResponse, error) {
	return uc.classroomIntegration.ListTeachers(courseID)
}
func (uc *ClassroomUseCase) DeleteTeacher(courseID, userID string) error {
	return uc.classroomIntegration.DeleteTeacher(courseID, userID)
}

// --- Invitations ---
func (uc *ClassroomUseCase) CreateInvitation(invitation *gclassroom.Invitation) (*gclassroom.Invitation, error) {
	return uc.classroomIntegration.CreateInvitation(invitation)
}
func (uc *ClassroomUseCase) GetInvitation(id string) (*gclassroom.Invitation, error) {
	return uc.classroomIntegration.GetInvitation(id)
}
func (uc *ClassroomUseCase) ListInvitations() (*gclassroom.ListInvitationsResponse, error) {
	return uc.classroomIntegration.ListInvitations()
}
func (uc *ClassroomUseCase) DeleteInvitation(id string) error {
	return uc.classroomIntegration.DeleteInvitation(id)
}
func (uc *ClassroomUseCase) AcceptInvitation(id string) error {
	return uc.classroomIntegration.AcceptInvitation(id)
}

// --- User Profiles ---
func (uc *ClassroomUseCase) GetUserProfile(userID string) (*gclassroom.UserProfile, error) {
	return uc.classroomIntegration.GetUserProfile(userID)
}

// --- Topics ---
func (uc *ClassroomUseCase) CreateTopic(courseID string, topic *gclassroom.Topic) (*gclassroom.Topic, error) {
	return uc.classroomIntegration.CreateTopic(courseID, topic)
}
func (uc *ClassroomUseCase) GetTopic(courseID, topicID string) (*gclassroom.Topic, error) {
	return uc.classroomIntegration.GetTopic(courseID, topicID)
}
func (uc *ClassroomUseCase) ListTopics(courseID string) (*gclassroom.ListTopicResponse, error) {
	return uc.classroomIntegration.ListTopics(courseID)
}
func (uc *ClassroomUseCase) UpdateTopic(courseID, topicID string, topic *gclassroom.Topic) (*gclassroom.Topic, error) {
	return uc.classroomIntegration.UpdateTopic(courseID, topicID, topic)
}
func (uc *ClassroomUseCase) DeleteTopic(courseID, topicID string) error {
	return uc.classroomIntegration.DeleteTopic(courseID, topicID)
}

// --- CourseWork ---
func (uc *ClassroomUseCase) CreateCourseWork(courseID string, courseWork *gclassroom.CourseWork) (*gclassroom.CourseWork, error) {
	return uc.classroomIntegration.CreateCourseWork(courseID, courseWork)
}
func (uc *ClassroomUseCase) GetCourseWork(courseID, courseWorkID string) (*gclassroom.CourseWork, error) {
	return uc.classroomIntegration.GetCourseWork(courseID, courseWorkID)
}
func (uc *ClassroomUseCase) ListCourseWork(courseID string) (*gclassroom.ListCourseWorkResponse, error) {
	return uc.classroomIntegration.ListCourseWork(courseID)
}
func (uc *ClassroomUseCase) UpdateCourseWork(courseID, courseWorkID string, courseWork *gclassroom.CourseWork) (*gclassroom.CourseWork, error) {
	return uc.classroomIntegration.UpdateCourseWork(courseID, courseWorkID, courseWork)
}
func (uc *ClassroomUseCase) DeleteCourseWork(courseID, courseWorkID string) error {
	return uc.classroomIntegration.DeleteCourseWork(courseID, courseWorkID)
}

// --- StudentSubmissions ---
func (uc *ClassroomUseCase) GetStudentSubmission(courseID, courseWorkID, submissionID string) (*gclassroom.StudentSubmission, error) {
	return uc.classroomIntegration.GetStudentSubmission(courseID, courseWorkID, submissionID)
}
func (uc *ClassroomUseCase) ListStudentSubmissions(courseID, courseWorkID string) (*gclassroom.ListStudentSubmissionsResponse, error) {
	return uc.classroomIntegration.ListStudentSubmissions(courseID, courseWorkID)
}
func (uc *ClassroomUseCase) PatchStudentSubmission(courseID, courseWorkID, submissionID string, submission *gclassroom.StudentSubmission, updateMask string) (*gclassroom.StudentSubmission, error) {
	return uc.classroomIntegration.PatchStudentSubmission(courseID, courseWorkID, submissionID, submission, updateMask)
}
func (uc *ClassroomUseCase) TurnInStudentSubmission(courseID, courseWorkID, submissionID string) error {
	return uc.classroomIntegration.TurnInStudentSubmission(courseID, courseWorkID, submissionID)
}
func (uc *ClassroomUseCase) ReturnStudentSubmission(courseID, courseWorkID, submissionID string) error {
	return uc.classroomIntegration.ReturnStudentSubmission(courseID, courseWorkID, submissionID)
}
func (uc *ClassroomUseCase) ReclaimStudentSubmission(courseID, courseWorkID, submissionID string) error {
	return uc.classroomIntegration.ReclaimStudentSubmission(courseID, courseWorkID, submissionID)
}

// --- Announcements ---
func (uc *ClassroomUseCase) CreateAnnouncement(courseID string, announcement *gclassroom.Announcement) (*gclassroom.Announcement, error) {
	return uc.classroomIntegration.CreateAnnouncement(courseID, announcement)
}
func (uc *ClassroomUseCase) GetAnnouncement(courseID, announcementID string) (*gclassroom.Announcement, error) {
	return uc.classroomIntegration.GetAnnouncement(courseID, announcementID)
}
func (uc *ClassroomUseCase) ListAnnouncements(courseID string) (*gclassroom.ListAnnouncementsResponse, error) {
	return uc.classroomIntegration.ListAnnouncements(courseID)
}
func (uc *ClassroomUseCase) UpdateAnnouncement(courseID, announcementID string, announcement *gclassroom.Announcement) (*gclassroom.Announcement, error) {
	return uc.classroomIntegration.UpdateAnnouncement(courseID, announcementID, announcement)
}
func (uc *ClassroomUseCase) DeleteAnnouncement(courseID, announcementID string) error {
	return uc.classroomIntegration.DeleteAnnouncement(courseID, announcementID)
}

// --- CourseWorkMaterials ---
func (uc *ClassroomUseCase) CreateCourseWorkMaterial(courseID string, material *gclassroom.CourseWorkMaterial) (*gclassroom.CourseWorkMaterial, error) {
	return uc.classroomIntegration.CreateCourseWorkMaterial(courseID, material)
}
func (uc *ClassroomUseCase) GetCourseWorkMaterial(courseID, materialID string) (*gclassroom.CourseWorkMaterial, error) {
	return uc.classroomIntegration.GetCourseWorkMaterial(courseID, materialID)
}
func (uc *ClassroomUseCase) ListCourseWorkMaterials(courseID string) (*gclassroom.ListCourseWorkMaterialResponse, error) {
	return uc.classroomIntegration.ListCourseWorkMaterials(courseID)
}
func (uc *ClassroomUseCase) UpdateCourseWorkMaterial(courseID, materialID string, material *gclassroom.CourseWorkMaterial) (*gclassroom.CourseWorkMaterial, error) {
	return uc.classroomIntegration.UpdateCourseWorkMaterial(courseID, materialID, material)
}
func (uc *ClassroomUseCase) DeleteCourseWorkMaterial(courseID, materialID string) error {
	return uc.classroomIntegration.DeleteCourseWorkMaterial(courseID, materialID)
}
