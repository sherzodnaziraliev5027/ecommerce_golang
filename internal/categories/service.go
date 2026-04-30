package categories

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}
type CategoryNode struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Children []CategoryNode `json:"children"`
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// 🔹 Create category
func (s *Service) Create(name string, parentID *uuid.UUID) error {
	category := &Category{
		ID:       uuid.New(),
		Name:     name,
		ParentID: parentID,
	}

	return s.repo.Create(category)
}

// 🔹 Get all categories
func (s *Service) GetAll() ([]Category, error) {
	return s.repo.FindAll()
}

func (s *Service) GetTree() ([]CategoryNode, error) {
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// map of id → node pointer
	nodeMap := make(map[string]*CategoryNode)

	// initialize all nodes
	for _, c := range categories {
		id := c.ID.String()
		nodeMap[id] = &CategoryNode{
			ID:       id,
			Name:     c.Name,
			Children: []CategoryNode{},
		}
	}

	var roots []*CategoryNode

	// build relationships
	for _, c := range categories {
		id := c.ID.String()
		node := nodeMap[id]

		if c.ParentID == nil {
			roots = append(roots, node)
		} else {
			parentID := c.ParentID.String()
			if parent, ok := nodeMap[parentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	// convert []*CategoryNode → []CategoryNode
	var result []CategoryNode
	for _, r := range roots {
		result = append(result, *r)
	}

	return result, nil
}
