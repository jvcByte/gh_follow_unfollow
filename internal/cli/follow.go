package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jvcByte/gh_follow_unfollow/internal/config"
	"github.com/jvcByte/gh_follow_unfollow/internal/git_hub_manager"
	"github.com/jvcByte/gh_follow_unfollow/internal/helper"
	"github.com/jvcByte/gh_follow_unfollow/internal/worker"
	"github.com/spf13/cobra"
)

func NewFollowCommand() *cobra.Command {
	cmd := followCommand()
	cmd.Flags().BoolP("force", "f", false, "Force follow")
	cmd.Flags().IntP("limit", "l", 0, "Max number of users to follow (0 = no limit)")

	return cmd
}

func followCommand() *cobra.Command {
	return &cobra.Command{ //nolint:gochecknoglobals // need for init command
		Use:   "follow <username>",
		Short: "Follow users",
		Long:  "Follow <username> followers who not follow you",
		Args:  cobra.ExactArgs(1), //nolint:mnd // args count
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err                            error
				force                          bool
				limit                          int
				cfg                            *config.Config
				username                       string
				followers, toFollow, following []string
			)
			username = args[0]
			cfg, err = config.Load()
			if err != nil {
				return err
			}
			force, err = cmd.Flags().GetBool("force")
			if err != nil {
				return fmt.Errorf("failed to get force flag: %w", err)
			}
			limit, err = cmd.Flags().GetInt("limit")
			if err != nil {
				return fmt.Errorf("failed to get limit flag: %w", err)
			}
			gm := git_hub_manager.NewGitHubManager(cfg.GitHubToken, cfg.GitHubUsername)
			followers, err = gm.GetFollowers(&username, limit)
			if err != nil {
				return err
			}
			following, err = gm.GetFollowing()
			if err != nil {
				return err
			}
			following = append(following, cfg.GitHubUsername)
			toFollow = gm.DiffUsernames(following, followers)

			if limit > 0 && len(toFollow) > limit {
				toFollow = toFollow[:limit]
			}

			if len(toFollow) == 0 {
				fmt.Println("No users to follow")
				return nil
			}

			for i, user := range toFollow {
				fmt.Printf("%d. %s\n", i+1, user)
			}

			if force == false {
				fmt.Printf("\nYou are sure that you want to follow %d users? (y/n):", len(toFollow))
				confirm := helper.GetInput("")
				if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
					fmt.Println("Operation canceled")
					return nil
				}
			}

			wp := worker.NewWorker(cfg.WorkerCount, cfg.QueueSize)
			wp.Start()

			for _, user := range toFollow {
				wp.AddTask(func() {
					err = gm.FollowUser(user, cfg.TimeDelay)
					if err != nil {
						fmt.Printf("Error follow %s: %v", user, err)
						os.Exit(1)
					}
					fmt.Printf("Followed %s\n", user)
				})

			}
			wp.Wait()
			wp.Stop()

			return nil
		},
	}
}