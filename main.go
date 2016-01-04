package main

import (
	"log"
	"net/http"
	"github.com/emicklei/go-restful"
	"errors"
	"os"
	"encoding/json"
	"github.com/daraszkrzysztof/secure-voting/logic"
)

type SecureVotingService struct {
	sv *logic.SecureVoting
}

func NewSecureVotingService() *SecureVotingService{
	return &SecureVotingService{sv : logic.NewSecureVoting()}
}

func (svs SecureVotingService)checkAdmin(request *restful.Request, response *restful.Response) (finish bool) {

	log.Printf("[secure-voting-service] check admin")

	_, err := svs.sv.CheckAdmin(request.PathParameter("admin-login"), request.HeaderParameter("organizer-password"))
	if err != nil {
		response.WriteError(http.StatusForbidden, err)
		return true
	}
	return false
}

func (svs SecureVotingService)findElection(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] find election")

	if svs.checkAdmin(request, response) {
		return
	}

	electionId := request.PathParameter("election-id")
	electionInfo, err := svs.sv.GetElectionInfo(electionId)

	if err != nil {
		response.WriteError(http.StatusNotFound, err)
	} else {
		response.WriteAsJson(electionInfo)
	}
}

func (svs SecureVotingService)listElections(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] list elections")

	if svs.checkAdmin(request, response) {
		return
	}

	response.WriteAsJson(svs.sv.ListElections)
}

func (svs SecureVotingService)createOrganiser(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] create organiser")

	orgName := request.PathParameter("organizer-name")
	orgPasswd := request.HeaderParameter("organizer-password")

	if org, err := svs.sv.CreateOrganizer(orgName, orgPasswd); err != nil {
		response.WriteError(http.StatusInternalServerError, errors.New("Not able to create new election."))
	}else {
		response.WriteAsJson(org)
	}

}

func (svs SecureVotingService)createElection(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] create election")

	err := request.Request.ParseForm()

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	orgName := request.PathParameter("organizer-name")
	orgPasswd := request.HeaderParameter("organizer-password")
	electionId := request.HeaderParameter("election-id")

	var initElConfig logic.Election

	decoder := json.NewDecoder(request.Request.Body)

	err = decoder.Decode(initElConfig)

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	svs.sv.CreateElections(orgName, orgPasswd, electionId, initElConfig)
}

func (svs SecureVotingService)createElectionBoardMembers(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] create election board members")

	err := request.Request.ParseForm()

	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (svs *SecureVotingService) Register() {

	ws := new(restful.WebService)

	ws.Path("/secure-voting").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	//admin interface

	ws.Route(ws.GET("/admin/{admin-login}/elections/{election-id}").
	To(svs.findElection).
	Param(ws.HeaderParameter("admin-password", "password for organizer's account")).
	Param(ws.PathParameter("admin-login", "admin login")).
	Param(ws.PathParameter("election-id", "id of the election")).
	Returns(200, "OK", []logic.Election{}))

	ws.Route(ws.GET("/admin/{admin-login}/list-elections").
	To(svs.listElections).
	Param(ws.HeaderParameter("organizer-password", "password for organizer's account")).
	Param(ws.PathParameter("admin-login", "admin login")).
	Returns(200, "OK", []string{}))

	//organizer interface

	ws.Route(ws.PUT("/organizer/{organizer-name}").To(svs.createOrganiser).
	Doc("Creates new organizer, if do not exist.").
	Operation("createOrganizer").
	Param(ws.HeaderParameter("organizer-password", "password for organizer's account")).
	Param(ws.PathParameter("organizer-name", "id of the organizer")).
	Writes(logic.Organizer{}))

	ws.Route(ws.POST("/{organizer-name}/election/{election-id}").To(svs.createElection).
	Doc("Creates new election for given organizer.").
	Operation("createNewElection").
	Param(ws.HeaderParameter("organizer-password", "password for organizer's account")).
	Param(ws.PathParameter("organizer-name", "id of the organizer")).
	Param(ws.PathParameter("election-id", "id of the new election")).
	Writes(logic.Organizer{}))

	ws.Route(ws.PUT("/{organizer-name}/election/{election-id}/board-members").To(svs.createElectionBoardMembers).
	Doc("Creates new election for given organizer.").
	Operation("createNewElection").
	Param(ws.HeaderParameter("organizer-password", "password for organizer's account")).
	Param(ws.PathParameter("organizer-name", "id of the organizer")).
	Param(ws.PathParameter("election-id", "id of the new election")).
	Writes(logic.Organizer{}))

	restful.Add(ws)
}

func main() {

	log.Printf("[secure-voting-service] init")

	svs := NewSecureVotingService()

	log.Printf("[secure-voting-service] register")

	svs.Register()

	log.Fatal(http.ListenAndServe(":" + os.Getenv("SECURE_VOTING_PORT"), nil))
}