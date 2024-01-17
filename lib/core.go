package server

/**
1. New Experiment is received
	- create Experiment
		- Experiment has between 2-4 options
		- Experiment may have 0-1 votes (if the creator already voted)
		- Experiment has an end date/time
2. Experiment is shared with someone
	- Creator enters a person's email address
	- Add the person as a Participant (?) on the Experiment
		- Experiments can't close until all Participants have completed
	- Send an email to the Participant
3. Participant votes
	- Participant uses link to create a session. 
		- There are no accounts (usernames/passwords) but all Experiments are forever associated with a set of Persons
	- Particpant chooses one of the 2-4 options
	- Participant can't see the outcome until Experiment closes
		- All Ps have voted, or E time has passed
4. Experiment ends
	- All Participants have voted
	- Experiment end time has passed
	- Calculate the outcome
		- If an option has received the most votes, it wins
		- In the event of a tie, flip a coin!
	- Notify everyone
*/

type Experiment struct {
	id int
	prompt string
	options []Option
	participants []Participant
	organiser Participant
}

// Returns the winning Option and the number of votes it received
func (e Experiment) winner() (Option, int) {
	max := 0
	var winner Option

	for i := 0; i < len(e.options); i++ {
		if e.options[i].votes > max {
			max = e.options[i].votes
			winner = e.options[i]
		}
	}

	return winner, max
}


// Returns true if the Experiment is waiting on 1 or more Participants to vote
func (e Experiment) isOpen() (bool) {
	totalVotes := 0

	for i := 0; i < len(e.options); i++ {
		totalVotes += e.options[i].votes
	}

	if totalVotes < len(e.participants) {
		return true
	}

	return false
}


type Option struct {
	value string
	votes int
}

type Participant struct {
	id int
	name string
	email string
}