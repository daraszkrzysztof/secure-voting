package logic
import (
	"github.com/daraszkrzysztof/secure-voting/secure"
	"errors"
)

type User struct {
	Id         string
	PasswdHash string
}

type Organizer struct {
	UserData         User
	HoldingElections []Election
}

type Voter struct {
	Id, Email string
}

type Election struct {
	Id, Title string
	Options   []string
	Voters    []Voter
}

type SecureVoting struct {
	ActiveElections map[string]Election
	Organizers      map[string]Organizer
}

func NewSecureVoting() *SecureVoting {
	return &SecureVoting{
		ActiveElections: make(map[string]Election),
		Organizers: make(map[string]Organizer)}
}

func (sv *SecureVoting)ListElections() []string {

	keys := make([]string, len(sv.ActiveElections))

	i := 0
	for k := range sv.ActiveElections {
		keys[i] = k
		i++
	}
	return keys
}

func (sv *SecureVoting)GetElectionInfo(electionId string) (Election, error) {

	foundElection, ok := sv.ActiveElections[electionId]

	if ok {
		return Election{Id:foundElection.Id, Title:foundElection.Title, Options:foundElection.Options},nil
	} else {
		return Election{}, errors.New("Couldn't find "+electionId+" election.")
	}
}

func (sv *SecureVoting)CreateOrganizer(id string, passwd string) (Organizer, error) {

	if _, ok := sv.Organizers[id]; ok {

		return Organizer{},errors.New("Organizer "+id+" already exists.")
	}else{
		passwd :=secure.DoUserPasswdHash(passwd)
		organizer := Organizer{User{Id:id, PasswdHash: passwd}, []Election{}}
		sv.Organizers[id] = organizer
	}
	return sv.Organizers[id],nil
}

func (sv *SecureVoting)checkOrganizer(organizerId, organizerPasswd string) (Organizer, error) {

	if foundOrganizer, okFound := sv.Organizers[organizerId]; okFound {
		if foundOrganizer.UserData.PasswdHash == secure.DoUserPasswdHash(organizerPasswd) {
			return foundOrganizer, nil
		}else{
			return Organizer{},errors.New("Incorrect password.")
		}
	}else{
		return Organizer{},errors.New("Organizer not found.")
	}
}

func (sv *SecureVoting)CreateElections(organizerId, organizerPasswd, electionId string, initElectionConfig Election) {

	if organizer, err := sv.checkOrganizer(organizerId, organizerPasswd); err!=nil {

		organizer.HoldingElections = append(organizer.HoldingElections, initElectionConfig)
		sv.ActiveElections[electionId] = initElectionConfig
	}
}

func (u *User)CreateElections(organizerId, title string) (Election, error) {
	return Election{Id: "1", Title: title, Options: []string{"Candidate A", "Candidate B"}}, nil
}