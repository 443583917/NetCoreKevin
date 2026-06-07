package entity

import "time"

type AIApp struct {
	ID           int64     `json:"id"`
	AppName      string    `json:"appName"`
	AppDesc      string    `json:"appDesc"`
	ModelID      int64     `json:"modelId"`
	SystemPrompt string    `json:"systemPrompt"`
	TenantID     int64     `json:"tenantId"`
	IsDelete     bool      `json:"isDelete"`
	CreateTime   time.Time `json:"createTime"`

	Model  *AIModel `json:"model,omitempty"`
	Skills []*Skill `json:"skills,omitempty" gorm:"many2many:t_ai_apps_bind_skill;"`
}

type AIModel struct {
	ID         int64     `json:"id"`
	ModelName  string    `json:"modelName"`
	Provider   string    `json:"provider"`
	APIKey     string    `json:"-"`
	BaseURL    string    `json:"baseUrl"`
	MaxTokens  int       `json:"maxTokens"`
	IsDelete   bool      `json:"isDelete"`
	CreateTime time.Time `json:"createTime"`
}

type Skill struct {
	ID          int64     `json:"id"`
	SkillName   string    `json:"skillName"`
	SkillCode   string    `json:"skillCode"`
	Description string    `json:"description"`
	SkillType   string    `json:"skillType"`
	Config      string    `json:"config"`
	IsDelete    bool      `json:"isDelete"`
	CreateTime  time.Time `json:"createTime"`
}

type ChatSession struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	AppID      int64     `json:"appId"`
	Title      string    `json:"title"`
	IsDelete   bool      `json:"isDelete"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type ChatMessage struct {
	ID         int64     `json:"id"`
	SessionID  int64     `json:"sessionId"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	Tokens     int       `json:"tokens"`
	CreateTime time.Time `json:"createTime"`
}

type KnowledgeBase struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	VectorModel string    `json:"vectorModel"`
	TenantID    int64     `json:"tenantId"`
	IsDelete    bool      `json:"isDelete"`
	CreateTime  time.Time `json:"createTime"`
}

type Document struct {
	ID         int64     `json:"id"`
	KBID       int64     `json:"kbId"`
	FileName   string    `json:"fileName"`
	FileURL    string    `json:"fileUrl"`
	Status     string    `json:"status"`
	ChunkCount int       `json:"chunkCount"`
	IsDelete   bool      `json:"isDelete"`
	CreateTime time.Time `json:"createTime"`
}
