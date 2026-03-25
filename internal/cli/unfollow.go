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

func NewUnFollowCommand() *cobra.Command {
	cmd := unFollowCommand()
	cmd.Flags().BoolP("force", "f", false, "Force unfollow")
	cmd.Flags().StringSlice("users", []string{}, "Specific users to unfollow if they don't follow back")

	return cmd
}

func unFollowCommand() *cobra.Command {
	return &cobra.Command{ //nolint:gochecknoglobals // need for init command
		Use:   "unfollow",
		Short: "Unfollow users",
		Long:  "Unfollow users who not follow you",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err                              error
				force                            bool
				targetUsers                      []string
				cfg                              *config.Config
				followers, toUnfollow, following []string
			)
			cfg, err = config.Load()
			if err != nil {
				return err
			}
			force, err = cmd.Flags().GetBool("force")
			if err != nil {
				return fmt.Errorf("failed to get force flag: %w", err)
			}
			targetUsers, err = cmd.Flags().GetStringSlice("users")
			if err != nil {
				return fmt.Errorf("failed to get users flag: %w", err)
			}
			gm := git_hub_manager.NewGitHubManager(cfg.GitHubToken, cfg.GitHubUsername)
			followers, err = gm.GetFollowers(nil, 0)
			if err != nil {
				return err
			}

			if len(targetUsers) > 0 {
				// Only consider the specified users, unfollow those not in followers
				following = targetUsers
			} else {
				following, err = gm.GetFollowing()
				if err != nil {
					return err
				}
			}

			toUnfollow = gm.DiffUsernames(followers, following)

			for i, user := range toUnfollow {
				fmt.Printf("%d. %s\n", i+1, user)
			}
			if len(toUnfollow) == 0 {
				fmt.Println("No users to unfollow")
				return nil
			}
			if force == false {
				fmt.Printf("\nYou are shure that you want to unfollow %d users?", len(toUnfollow))
				confirm := helper.GetInput("")
				if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
					fmt.Println("Operation canceled")
					return nil
				}
			}

			wp := worker.NewWorker(cfg.WorkerCount, cfg.QueueSize)
			wp.Start()

			for _, user := range toUnfollow {
				wp.AddTask(func() {
					err = gm.UnfollowUser(user, cfg.TimeDelay)
					if err != nil {
						fmt.Printf("Error unfollow %s: %v", user, err)
						os.Exit(1)
					}
					fmt.Printf("UnFollowed %s\n", user)
				})

			}
			wp.Wait()
			wp.Stop()

			return nil
		},
	}
}
