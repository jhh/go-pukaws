package storage // import "jhhgo.us/pukaws/storage"

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"jhhgo.us/pukaws/model"
)

// NewBookmarkMgoStorage initializes the storage
func NewBookmarkMgoStorage(uri string) (BookmarkMgoStorage, error) {
	session, err := mgo.Dial(uri)
	if err != nil {
		return BookmarkMgoStorage{}, err
	}
	return BookmarkMgoStorage{session: session}, nil
}

// BookmarkMgoStorage stores all users
type BookmarkMgoStorage struct {
	session *mgo.Session
}

// Close will close the master session
func (s BookmarkMgoStorage) Close() {
	s.session.Close()
}

// GetAll returns the bookmarks specified by query.
func (s BookmarkMgoStorage) GetAll(q Query) ([]lib.Bookmark, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	iter := c.Find(q.Mgo()).Sort("-timestamp").Iter()
	var result []lib.Bookmark
	err := iter.All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetPage returns a portion of bookmarks specified by query.
func (s BookmarkMgoStorage) GetPage(q Query, skip, limit int) ([]lib.Bookmark, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	iter := c.Find(q.Mgo()).Sort("-timestamp").Skip(skip).Limit(limit).Iter()
	var result []lib.Bookmark
	err := iter.All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Count returns the total number of bookmarks specified by query.
func (s BookmarkMgoStorage) Count(q Query) (int, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	n, err := c.Find(q.Mgo()).Count()
	return n, err
}

// GetOne user
func (s BookmarkMgoStorage) GetOne(id string) (lib.Bookmark, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	var result lib.Bookmark
	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return lib.Bookmark{}, err
	}
	return result, nil
}

// Insert a user and set Timestamp to insert time if not already set.
func (s BookmarkMgoStorage) Insert(b *lib.Bookmark) error {
	id := bson.NewObjectId()
	b.ID = id

	if b.Timestamp.IsZero() {
		b.Timestamp = time.Now()
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	err := c.Insert(b)
	return err
}

// Delete one :(
func (s BookmarkMgoStorage) Delete(id string) error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	err := c.RemoveId(bson.ObjectIdHex(id))
	return err
}

// Update a user and updates Timestamp.
func (s BookmarkMgoStorage) Update(b *lib.Bookmark) error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C("bookmarks")
	b.Timestamp = time.Now()
	err := c.UpdateId(b.ID, b)
	return err
}