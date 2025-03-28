package crowdfund

import (
	"encoding/json"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) createHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	activeCampaign := c.activeCampaign()
	if activeCampaign != nil {
		return cmd.RenderFailedTemplateF("There is an active campaign: %s", activeCampaign.Title)
	}

	title := args[argNameCreateTitle]
	desc := args[argNameCreateDesc]
	packagesJSON := args[argNameCreatePackages]

	packages := []entity.Package{}
	err := json.Unmarshal([]byte(packagesJSON), &packages)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	if title == "" {
		return cmd.RenderFailedTemplate("The title of the crowdfunding campaign cannot be empty")
	}

	campaign := &entity.CrowdfundCampaign{
		CreatorID: caller.ID,
		Title:     title,
		Desc:      desc,
		Packages:  packages,
		Active:    true,
	}
	err = c.db.AddCrowdfundCampaign(campaign)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("campaign", campaign)
}
