package Structs


type UserInfo struct {
	IdPart  string  
	IdGrp  	int32   
	IdUsr  	int32   
	Nombre 	string  
	Status 	bool    
	PathD	string	 
}

var UsuarioActual UserInfo

func SalirUsuario(){
	UsuarioActual.IdGrp = 0
	UsuarioActual.IdPart = ""
	UsuarioActual.IdUsr = 0
	UsuarioActual.Nombre = ""
	UsuarioActual.Status = false
	UsuarioActual.PathD = ""
}


 

 
 
 