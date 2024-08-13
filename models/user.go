package models

type User struct {
	//gorm.Model
	ID                uint   `gorm:"primaryKey;column:cod_persona" json:"cod_persona"`
	Nombre            string `json:"Nombre"`
	Apellidos         string `json:"Apellidos"`
	Rut               string `json:"Rut"`
	Nacionalidad      string `json:"nacionalidad"`
	MailSid           string `json:"mail_sid"`
	MailInstitucional string `json:"mail_institucional"`
	MailOpcional      string `json:"mail_opcional"`
	NombreCompleto    string `json:"nombre_completo"`
	Codsexo           uint   `json:"cod_sexo"`
	Dessexo           string `json:"des_sexo"`
	FechaNac          string `json:"fecha_nac"`
	Ap1Persona        string `json:"ap1_persona"`
	Ap2Persona        string `json:"ap2_persona"`
	EsChileno         string `json:"eschileno"`
	Roles             []Role `gorm:"many2many:detalle_categoria;foreignKey:ID;joinForeignKey:CodPersona;References:ID;joinReferences:CodCategoria"`
}

func (User) TableName() string {
	return "vista_personas"
}

type Role struct {
	//gorm.Model
	ID          uint   `gorm:"primaryKey;column:cod_categoria" json:"cod_categoria"`
	Descripcion string `json:"descripcion"`
	IdAcceso    uint   `json:"id_acceso"`
}

func (Role) TableName() string {
	return "categorias"
}
