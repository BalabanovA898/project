package usecase

import (
	"context"

	"github.com/BalabanovA898/project/account/src/domain"
	"github.com/BalabanovA898/project/account/src/repository/postgres"
)

type TeamUseCase struct {
	teams *postgres.TeamRepository
}

func NewTeamUseCase(teams *postgres.TeamRepository) *TeamUseCase {
	return &TeamUseCase{teams: teams}
}

func (u *TeamUseCase) CreateTeam(ctx context.Context, name, ownerID string) (domain.Team, error) {
	team := domain.Team{
		Name:    name,
		OwnerID: ownerID,
	}

	id, err := u.teams.Create(ctx, team)
	if err != nil {
		return domain.Team{}, err
	}

	team.ID = id

	// Получаем команду с членами
	return u.teams.GetByID(ctx, id)
}

func (u *TeamUseCase) GetTeam(ctx context.Context, teamID string) (domain.Team, error) {
	return u.teams.GetByID(ctx, teamID)
}

func (u *TeamUseCase) JoinTeam(ctx context.Context, teamID, userID string) error {
	return u.teams.AddMember(ctx, teamID, userID)
}

func (u *TeamUseCase) LeaveTeam(ctx context.Context, teamID, userID string) error {
	return u.teams.RemoveMember(ctx, teamID, userID)
}
