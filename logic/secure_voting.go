package logic
import (
	"github.com/daraszkrzysztof/secure-voting/secure"
	"errors"
)

type Organizer struct {
	Id         		string
	PasswdHash 		string
	HoldingElections map[string]Election
}

/**
the most secure way would be :
1. send an email with confirmation
2. after clicking the link in the email, he would received PrivateKey in base64 format
3. after sending private key, public key is saved in a structure

the simplest solution : define board member, in response receive private key
 */
type BoardMember struct {
	Name,Email		string
	PublicKeyDer 	[]byte
}

type Voter struct {
	Name, Email 	string
}

type Election struct {
	Id, Title 	string
	Options   	[]string
	Voters    	[]Voter
	BoardMember []BoardMember
}

type SecureVoting struct {
	Admins			map[string]Admin
	ActiveElections map[string]Election
	Organizers      map[string]Organizer
}

func NewSecureVoting() *SecureVoting {
	return &SecureVoting{
		ActiveElections: make(map[string]Election),
		Organizers: make(map[string]Organizer)}
}

func (sv *SecureVoting)CheckAdmin(adminId, adminPasswd string) (Admin, error){

	admin := sv.Admins[adminId]
	if &admin!=nil && admin.authenticate(adminPasswd){
		return admin, nil
	}else {
		return Admin{},errors.New("Admin "+adminId+" does not exist or passwd invalid.")
	}
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
		opasswd :=secure.DoUserPasswdHash(passwd)
		organizer := Organizer{ id, opasswd, make(map[string]Election) }
		sv.Organizers[id] = organizer
	}
	return sv.Organizers[id],nil
}

func (sv *SecureVoting)checkOrganizer(organizerId, organizerPasswd string) (Organizer, error) {

	if foundOrganizer, okFound := sv.Organizers[organizerId]; okFound {
		if foundOrganizer.PasswdHash == secure.DoUserPasswdHash(organizerPasswd) {
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

		organizer.HoldingElections[electionId] = initElectionConfig
		sv.ActiveElections[electionId] = initElectionConfig
	}
}

func (u *Organizer)CreateElections(organizerId, title string) (Election, error) {
	return Election{Id: "1", Title: title, Options: []string{"Candidate A", "Candidate B"}}, nil
}