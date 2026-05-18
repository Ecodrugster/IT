package models

import "time"

type Post struct {
	ID        string    `json:"id" firestore:"-"`
	AuthorID  string    `json:"author_id" firestore:"author_id"`
	Content   string    `json:"content" firestore:"content"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	Likes     int       `json:"likes" firestore:"likes"`
}

type Comment struct {
	ID        string    `json:"id" firestore:"-"`
	PostID    string    `json:"post_id" firestore:"post_id"`
	AuthorID  string    `json:"author_id" firestore:"author_id"`
	Text      string    `json:"text" firestore:"text"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

type User struct {
	UID         string `json:"uid" firestore:"uid"`
	DisplayName string `json:"display_name" firestore:"display_name"`
	Email       string `json:"email" firestore:"email"`
	PhotoURL    string `json:"photo_url" firestore:"photo_url"`
	Role        string `json:"role" firestore:"role"` // 'student' or 'admin'
}
