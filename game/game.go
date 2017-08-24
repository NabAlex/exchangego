package game

import (
	"errors"
	"sync"
	"time"

	"github.com/shmel1k/exchangego/database"
	"github.com/shmel1k/exchangego/exchange"
)

type game struct {
	duration int64
	end      int64

	move bool // False -- down, True -- up
}

const (
	updateTime = 1 * time.Second
)

var (
	ErrUserExists = errors.New("failed to add user to game: user exists")
)

var players Players

type Players struct {
	players map[exchange.User]game

	mu sync.Mutex
}

func (p *Players) Add(user exchange.User, duration int64, move bool) error {
	if p.players != nil {
		if _, ok := p.players[user]; ok {
			return ErrUserExists
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.players == nil {
		p.players = make(map[exchange.User]game)
	}

	p.players[user] = game{
		duration: duration,
		end:      time.Now().Unix() + duration,
		move:     move,
	}

	return nil
}

func (p *Players) Delete(user exchange.User) {
	p.mu.Lock()
	delete(p.players, user)
	p.mu.Unlock()
}

func AddPlayer(user exchange.User, duration int64, move bool) error {
	return players.Add(user, duration, move)
}

func RunScheduler() error {
	for {
		err := schedule()
		if err != nil {
			return err
		}
	}
}

func schedule() error {
	playersToUpdate := make([]exchange.User, 0, len(players.players))
	for {
		t := time.Now().Unix()
		for k, v := range players.players {
			if v.end <= t {
				playersToUpdate = append(playersToUpdate, k)
			}
		}
		var err error
		for _, v := range playersToUpdate {
			err = database.UpdateMoney(v.ID, v.Money/2)
			if err != nil {
				return err
			}
		}
		playersToUpdate = playersToUpdate[:0]

		time.Sleep(updateTime)
	}
}
