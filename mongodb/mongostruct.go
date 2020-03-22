package mongodb

type BlogItem struct {
	ID      string `bson:"_id,omitempty"`
	AutorID string `bson:"author_id"`
	Content string `bson:"content"`
	Title   string `bson:"title"`
}
