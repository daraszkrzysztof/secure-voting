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

func (svs SecureVotingService)findElection(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] find election")

	electionId := request.PathParameter("election-id")
	electionInfo, err := svs.sv.GetElectionInfo(electionId)

	if err != nil {
		response.WriteAsJson(electionInfo)
	} else {
		response.WriteError(http.StatusNotFound, err)
	}
}

func (svs SecureVotingService)listElections(request *restful.Request, response *restful.Response) {

	log.Printf("[secure-voting-service] list elections")

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

func (svs *SecureVotingService) Register() {

	ws := new(restful.WebService)

	ws.Path("/secure-voting").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/elections/{election-id}").
	To(svs.findElection).
	Param(ws.PathParameter("election-id", "id of the election")).
	Returns(200, "OK", []logic.Election{}))

	ws.Route(ws.GET("/list-elections").
	To(svs.listElections).
	Returns(200, "OK", []string{}))

	ws.Route(ws.PUT("/new-organizer/{organizer-name}").To(svs.createOrganiser).
	Doc("Adds new organizer.").
	Operation("createNewOrganizer").
	Param(ws.HeaderParameter("organizer-password", "password of the organizer")).
	Param(ws.PathParameter("organizer-name", "id of the organizer")).
	Writes(logic.Organizer{}))

	ws.Route(ws.POST("/{organizer-name}/new-election/{election-id}").To(svs.createElection).
	Doc("Creates new election for given organizer.").
	Operation("createNewElection").
	Param(ws.HeaderParameter("organizer-password", "password of the organizer")).
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