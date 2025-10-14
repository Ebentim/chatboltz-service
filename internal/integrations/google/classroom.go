package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

// ClassroomService provides methods for interacting with the Google Classroom API.
type ClassroomService struct {
	service *classroom.Service
}

// NewClassroomService creates a new ClassroomService.
func NewClassroomService(ctx context.Context, ts oauth2.TokenSource) (*ClassroomService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := classroom.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Classroom client: %v", err)
	}
	return &ClassroomService{service: srv}, nil
}

// --- Courses ---

// CreateCourse creates a new course.
func (s *ClassroomService) CreateCourse(course *classroom.Course) (*classroom.Course, error) {
	return s.service.Courses.Create(course).Do()
}

// GetCourse retrieves a course.
func (s *ClassroomService) GetCourse(courseID string) (*classroom.Course, error) {
	return s.service.Courses.Get(courseID).Do()
}

// ListCourses lists courses.
func (s *ClassroomService) ListCourses() (*classroom.ListCoursesResponse, error) {
	return s.service.Courses.List().Do()
}

// UpdateCourse updates a course.
func (s *ClassroomService) UpdateCourse(courseID string, course *classroom.Course) (*classroom.Course, error) {
	return s.service.Courses.Update(courseID, course).Do()
}

// DeleteCourse deletes a course.
func (s *ClassroomService) DeleteCourse(courseID string) error {
	_, err := s.service.Courses.Delete(courseID).Do()
	return err
}

// --- Aliases ---

// CreateAlias creates a new alias for a course.
func (s *ClassroomService) CreateAlias(courseID string, alias *classroom.CourseAlias) (*classroom.CourseAlias, error) {
	return s.service.Courses.Aliases.Create(courseID, alias).Do()
}

// ListAliases lists the aliases for a course.
func (s *ClassroomService) ListAliases(courseID string) (*classroom.ListCourseAliasesResponse, error) {
	return s.service.Courses.Aliases.List(courseID).Do()
}

// DeleteAlias deletes an alias for a course.
func (s *ClassroomService) DeleteAlias(courseID, alias string) error {
	_, err := s.service.Courses.Aliases.Delete(courseID, alias).Do()
	return err
}

// --- Students ---

// CreateStudent adds a student to a course.
func (s *ClassroomService) CreateStudent(courseID string, student *classroom.Student) (*classroom.Student, error) {
	return s.service.Courses.Students.Create(courseID, student).Do()
}

// GetStudent retrieves a student from a course.
func (s *ClassroomService) GetStudent(courseID, userID string) (*classroom.Student, error) {
	return s.service.Courses.Students.Get(courseID, userID).Do()
}

// ListStudents lists the students in a course.
func (s *ClassroomService) ListStudents(courseID string) (*classroom.ListStudentsResponse, error) {
	return s.service.Courses.Students.List(courseID).Do()
}

// DeleteStudent removes a student from a course.
func (s *ClassroomService) DeleteStudent(courseID, userID string) error {
	_, err := s.service.Courses.Students.Delete(courseID, userID).Do()
	return err
}

// --- Teachers ---

// CreateTeacher adds a teacher to a course.
func (s *ClassroomService) CreateTeacher(courseID string, teacher *classroom.Teacher) (*classroom.Teacher, error) {
	return s.service.Courses.Teachers.Create(courseID, teacher).Do()
}

// GetTeacher retrieves a teacher from a course.
func (s *ClassroomService) GetTeacher(courseID, userID string) (*classroom.Teacher, error) {
	return s.service.Courses.Teachers.Get(courseID, userID).Do()
}

// ListTeachers lists the teachers in a course.
func (s *ClassroomService) ListTeachers(courseID string) (*classroom.ListTeachersResponse, error) {
	return s.service.Courses.Teachers.List(courseID).Do()
}

// DeleteTeacher removes a teacher from a course.
func (s *ClassroomService) DeleteTeacher(courseID, userID string) error {
	_, err := s.service.Courses.Teachers.Delete(courseID, userID).Do()
	return err
}

// --- Invitations ---

// CreateInvitation creates an invitation.
func (s *ClassroomService) CreateInvitation(invitation *classroom.Invitation) (*classroom.Invitation, error) {
	return s.service.Invitations.Create(invitation).Do()
}

// GetInvitation retrieves an invitation.
func (s *ClassroomService) GetInvitation(id string) (*classroom.Invitation, error) {
	return s.service.Invitations.Get(id).Do()
}

// ListInvitations lists invitations.
func (s *ClassroomService) ListInvitations() (*classroom.ListInvitationsResponse, error) {
	return s.service.Invitations.List().Do()
}

// DeleteInvitation deletes an invitation.
func (s *ClassroomService) DeleteInvitation(id string) error {
	_, err := s.service.Invitations.Delete(id).Do()
	return err
}

// AcceptInvitation accepts an invitation.
func (s *ClassroomService) AcceptInvitation(id string) error {
	_, err := s.service.Invitations.Accept(id).Do()
	return err
}

// --- User Profiles ---

// GetUserProfile retrieves a user's profile.
func (s *ClassroomService) GetUserProfile(userID string) (*classroom.UserProfile, error) {
	return s.service.UserProfiles.Get(userID).Do()
}

// --- Topics ---

// CreateTopic creates a topic.
func (s *ClassroomService) CreateTopic(courseID string, topic *classroom.Topic) (*classroom.Topic, error) {
	return s.service.Courses.Topics.Create(courseID, topic).Do()
}

// GetTopic retrieves a topic.
func (s *ClassroomService) GetTopic(courseID, topicID string) (*classroom.Topic, error) {
	return s.service.Courses.Topics.Get(courseID, topicID).Do()
}

// ListTopics lists topics in a course.
func (s *ClassroomService) ListTopics(courseID string) (*classroom.ListTopicResponse, error) {
	return s.service.Courses.Topics.List(courseID).Do()
}

// UpdateTopic updates a topic.
func (s *ClassroomService) UpdateTopic(courseID, topicID string, topic *classroom.Topic) (*classroom.Topic, error) {
	return s.service.Courses.Topics.Patch(courseID, topicID, topic).UpdateMask("name").Do()
}

// DeleteTopic deletes a topic.
func (s *ClassroomService) DeleteTopic(courseID, topicID string) error {
	_, err := s.service.Courses.Topics.Delete(courseID, topicID).Do()
	return err
}

// --- Course Work ---

// CreateCourseWork creates course work.
func (s *ClassroomService) CreateCourseWork(courseID string, courseWork *classroom.CourseWork) (*classroom.CourseWork, error) {
	return s.service.Courses.CourseWork.Create(courseID, courseWork).Do()
}

// GetCourseWork retrieves course work.
func (s *ClassroomService) GetCourseWork(courseID, courseWorkID string) (*classroom.CourseWork, error) {
	return s.service.Courses.CourseWork.Get(courseID, courseWorkID).Do()
}

// ListCourseWork lists course work for a course.
func (s *ClassroomService) ListCourseWork(courseID string) (*classroom.ListCourseWorkResponse, error) {
	return s.service.Courses.CourseWork.List(courseID).Do()
}

// UpdateCourseWork updates course work.
func (s *ClassroomService) UpdateCourseWork(courseID, courseWorkID string, courseWork *classroom.CourseWork) (*classroom.CourseWork, error) {
	return s.service.Courses.CourseWork.Patch(courseID, courseWorkID, courseWork).UpdateMask("title,description,dueDate,maxPoints").Do()
}

// DeleteCourseWork deletes course work.
func (s *ClassroomService) DeleteCourseWork(courseID, courseWorkID string) error {
	_, err := s.service.Courses.CourseWork.Delete(courseID, courseWorkID).Do()
	return err
}

// --- Student Submissions ---

// GetStudentSubmission retrieves a student submission.
func (s *ClassroomService) GetStudentSubmission(courseID, courseWorkID, submissionID string) (*classroom.StudentSubmission, error) {
	return s.service.Courses.CourseWork.StudentSubmissions.Get(courseID, courseWorkID, submissionID).Do()
}

// ListStudentSubmissions lists student submissions for a course work.
func (s *ClassroomService) ListStudentSubmissions(courseID, courseWorkID string) (*classroom.ListStudentSubmissionsResponse, error) {
	return s.service.Courses.CourseWork.StudentSubmissions.List(courseID, courseWorkID).Do()
}

// PatchStudentSubmission updates a student submission.
func (s *ClassroomService) PatchStudentSubmission(courseID, courseWorkID, submissionID string, submission *classroom.StudentSubmission, updateMask string) (*classroom.StudentSubmission, error) {
	return s.service.Courses.CourseWork.StudentSubmissions.Patch(courseID, courseWorkID, submissionID, submission).UpdateMask(updateMask).Do()
}

// TurnInStudentSubmission turns in a student submission.
func (s *ClassroomService) TurnInStudentSubmission(courseID, courseWorkID, submissionID string) error {
	_, err := s.service.Courses.CourseWork.StudentSubmissions.TurnIn(courseID, courseWorkID, submissionID, &classroom.TurnInStudentSubmissionRequest{}).Do()
	return err
}

// ReturnStudentSubmission returns a student submission.
func (s *ClassroomService) ReturnStudentSubmission(courseID, courseWorkID, submissionID string) error {
	_, err := s.service.Courses.CourseWork.StudentSubmissions.Return(courseID, courseWorkID, submissionID, &classroom.ReturnStudentSubmissionRequest{}).Do()
	return err
}

// ReclaimStudentSubmission reclaims a student submission.
func (s *ClassroomService) ReclaimStudentSubmission(courseID, courseWorkID, submissionID string) error {
	_, err := s.service.Courses.CourseWork.StudentSubmissions.Reclaim(courseID, courseWorkID, submissionID, &classroom.ReclaimStudentSubmissionRequest{}).Do()
	return err
}

// --- Announcements ---

// CreateAnnouncement creates an announcement.
func (s *ClassroomService) CreateAnnouncement(courseID string, announcement *classroom.Announcement) (*classroom.Announcement, error) {
	return s.service.Courses.Announcements.Create(courseID, announcement).Do()
}

// GetAnnouncement retrieves an announcement.
func (s *ClassroomService) GetAnnouncement(courseID, announcementID string) (*classroom.Announcement, error) {
	return s.service.Courses.Announcements.Get(courseID, announcementID).Do()
}

// ListAnnouncements lists announcements for a course.
func (s *ClassroomService) ListAnnouncements(courseID string) (*classroom.ListAnnouncementsResponse, error) {
	return s.service.Courses.Announcements.List(courseID).Do()
}

// UpdateAnnouncement updates an announcement.
func (s *ClassroomService) UpdateAnnouncement(courseID, announcementID string, announcement *classroom.Announcement) (*classroom.Announcement, error) {
	return s.service.Courses.Announcements.Patch(courseID, announcementID, announcement).UpdateMask("text").Do()
}

// DeleteAnnouncement deletes an announcement.
func (s *ClassroomService) DeleteAnnouncement(courseID, announcementID string) error {
	_, err := s.service.Courses.Announcements.Delete(courseID, announcementID).Do()
	return err
}

// --- CourseWorkMaterials ---

// CreateCourseWorkMaterial creates a course work material.
func (s *ClassroomService) CreateCourseWorkMaterial(courseID string, material *classroom.CourseWorkMaterial) (*classroom.CourseWorkMaterial, error) {
	return s.service.Courses.CourseWorkMaterials.Create(courseID, material).Do()
}

// GetCourseWorkMaterial retrieves a course work material.
func (s *ClassroomService) GetCourseWorkMaterial(courseID, materialID string) (*classroom.CourseWorkMaterial, error) {
	return s.service.Courses.CourseWorkMaterials.Get(courseID, materialID).Do()
}

// ListCourseWorkMaterials lists course work materials for a course.
func (s *ClassroomService) ListCourseWorkMaterials(courseID string) (*classroom.ListCourseWorkMaterialResponse, error) {
	return s.service.Courses.CourseWorkMaterials.List(courseID).Do()
}

// UpdateCourseWorkMaterial updates a course work material.
func (s *ClassroomService) UpdateCourseWorkMaterial(courseID, materialID string, material *classroom.CourseWorkMaterial) (*classroom.CourseWorkMaterial, error) {
	return s.service.Courses.CourseWorkMaterials.Patch(courseID, materialID, material).UpdateMask("title,description").Do()
}

// DeleteCourseWorkMaterial deletes a course work material.
func (s *ClassroomService) DeleteCourseWorkMaterial(courseID, materialID string) error {
	_, err := s.service.Courses.CourseWorkMaterials.Delete(courseID, materialID).Do()
	return err
}
