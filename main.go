package main

import (
	"log"
	"crypto/sha512"
	"net/http"
	"encoding/base64"
	"github.com/emicklei/go-restful"
	"errors"
)

type User struct{
	Id string
	PasswdHash string
}

type Organizer struct{
	UserData User
}

type Election struct{
	Id,Title string
	Options []string
}

type SecureVoting struct{
	ActiveElections []Election
	Organizers []Organizer
}

type SecureVotingService struct{

	sv SecureVoting
}

func (sv *SecureVoting)createNewOrganizer(id string, passwd string)(Organizer, error){

	hasher := sha512.New()
	hasher.Write( []byte(passwd) )
	passwdHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return Organizer{User{Id:id, PasswdHash: passwdHash}},nil
}

func (u *User)createNewElections(organizerId, title string) (Election, error){
	return Election{Id: "1", Title: title, Options: []string{"Candidate A", "Candidate B"}}, nil
}

func (svs SecureVotingService)findAllElections(request *restful.Request, response *restful.Response){

	log.Printf("[find all elections]")

	if &svs.sv.ActiveElections != nil {
		response.WriteAsJson(svs.sv.ActiveElections)
	} else {
		response.WriteError(http.StatusNotFound, errors.New("No elections found."))
	}

}

func (svs SecureVotingService)addNewOrganiser(request *restful.Request, response *restful.Response){

	log.Printf("[add new organiser]")

	orgName :=  request.PathParameter("organizer-name")
	orgPasswd := request.HeaderParameter("organizer-password")

	if org, err := svs.sv.createNewOrganizer(orgName, orgPasswd); err!=nil {
		response.WriteError(http.StatusInternalServerError, errors.New("Not able to create new election."))
	}else{
		response.WriteAsJson(org)
	}

}

func (svs SecureVotingService) Register()  {

	ws := new(restful.WebService)

	ws.
		Path("/secure-voting").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/all-elections").
		To(svs.findAllElections).
		Returns(200, "OK", []Election{}))

	ws.Route(ws.PUT("/new-organizer/{organizer-name}").To(svs.addNewOrganiser).
		Doc("Adds new organizer.").
		Operation("createNewOrganizer").
		Param(ws.HeaderParameter("organizer-password", "password of the organizer")).
		Param(ws.PathParameter("organizer-name", "id of the organizer")).
		Writes(Organizer{}))

	restful.Add(ws)
}

func main(){

	svs := SecureVotingService{}
	svs.Register()

/*	config := swagger.Config{
		WebServices: restful.RegisteredWebServices(),
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/home/krzysztof/dev/github/swagger-ui/dist",
	}
	swagger.InstallSwaggerService(config)*/

	log.Printf("[secure-voting] start.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}