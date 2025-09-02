package templates

import (
	"fmt"

	"circles.diy/internal/models"
	"circles.diy/internal/utils"
)

func GetMockDashboardData() models.DashboardData {
	return models.DashboardData{
		BaseData: models.BaseData{
			Title:     "Dashboard",
			ActiveNav: "dashboard",
			Theme: models.ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: utils.GenerateCSRFToken(),
		},
		Feed: []models.FeedItem{
			{
				ID: "1",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=48&h=48&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2m ago",
				Circle:  "Woodworking",
				Image: &models.MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=600&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8b2FrJTIwdGFibGV8ZW58MHx8MHx8fDI%3D",
					Alt: "Oak coffee table project",
				},
				CanBuy: true,
			},
			{
				ID: "2",
				User: models.User{
					ID:     "heathtyler",
					Handle: "@heathtyler",
					Avatar: "https://images.unsplash.com/photo-1581391528803-54be77ce23e3?w=48&h=48&fit=crop&crop=face",
				},
				Content: "Warming up the barbeque and just got couple cases of the finest bread-water. Keen to see you all... remember 7PM dont be late!",
				TimeAgo: "15m ago",
				Circle:  "The Crop Circle",
				Image: &models.MediaItem{
					URL: "https://images.unsplash.com/photo-1664463758574-e640a7a998d4?q=80&w=600&auto=format&fit=crop",
					Alt: "BBQ gathering setup",
				},
				CanBuy: false,
			},
			{
				ID: "3",
				User: models.User{
					ID:     "sara_pcb",
					Handle: "@sara_pcb",
					Avatar: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=48&h=48&auto=format&fit=crop",
				},
				Content: "New tutorial series starting: \"Arduino for Beginners\". First session covers basic circuits and programming fundamentals.",
				TimeAgo: "1h ago",
				Circle:  "DIY Electronics",
				Image: &models.MediaItem{
					URL: "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=600&fit=crop&crop=center",
					Alt: "Arduino tutorial setup",
				},
				CanBuy: false,
			},
			{
				ID: "4",
				User: models.User{
					ID:     "zucc",
					Handle: "@zucc",
					Avatar: "https://media.tenor.com/y1mYLo66EuoAAAAM/zucky.gifs",
				},
				Content: "circles.diy has changed the game forever. \n\nFuck, I wish I'd thought of that.",
				TimeAgo: "1h ago",
				Circle:  "Communication Software",
			},
		},
		FeedOffset: 4,
		Circles: []models.Circle{
			{ID: "1", Name: "Woodworking", Thumbnail: "https://images.unsplash.com/photo-1504148455328-c376907d081c?w=48&h=48&fit=crop&crop=center", MemberCount: "8 online"},
			{ID: "2", Name: "The Crop Circle", Thumbnail: "https://images.unsplash.com/photo-1611402858501-d3de70f3c67e?q=80&w=48&auto=format&fit=crop", MemberCount: "12 members"},
			{ID: "3", Name: "DIY Electronics", Thumbnail: "https://images.unsplash.com/photo-1518611012118-696072aa579a?w=48&h=48&fit=crop&crop=center", MemberCount: "24 members"},
		},
		Discussions: []models.Discussion{
			{ID: "1", Title: "Sustainable Materials: Where to Source?", Preview: "Looking for suppliers of ethically sourced hardwoods. What are your go-to sources for...", Circle: "Woodworking", TimeAgo: "23 replies • 45m ago"},
			{ID: "2", Title: "Pricing Creative Work: Community Wisdom", Preview: "How do you approach pricing custom commissions? Struggling to find the balance between...", Circle: "Local Artists", TimeAgo: "8 replies • 2h ago"},
		},
		Events: []models.Event{
			{ID: "1", Title: "Workshop: Intro to Drum & Bass Mixing", Time: "2:00 PM", Day: "31", Month: "Aug"},
			{ID: "2", Title: "Monthly Showcase", Time: "7:00 PM", Day: "02", Month: "Sep"},
		},
		Ripples: []models.Ripple{
			{ID: "1", User: "@maia", Content: "Quick progress update on the oak table project...", ExpiresIn: "1h left"},
			{ID: "2", User: "@craft_collective", Content: "Tip: pre-drill hardwood to avoid splitting", ExpiresIn: "18h left"},
		},
		MarketplaceItems: []models.MarketplaceItem{
			{ID: "1", Title: "Cordless Drill", Price: "$85", Image: "https://images.unsplash.com/photo-1504148455328-c376907d081c?w=128&fit=crop&crop=center", Location: "Sydney, NSW", TimeAgo: "3h ago"},
			{ID: "2", Title: "Ceramic Wheel & Tools", Price: "Trade", Image: "https://unsplash.com/photos/ZSgWcW70cTs/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8NXx8Y2VyYW1pYyUyMHdoZWVsfGVufDB8fHx8MTc1NjU1NjAxNXwy&force=true&w=128", Location: "Stanwell Park, NSW", TimeAgo: "1d ago"},
			{ID: "3", Title: "Oak Lumber Bundle", Price: "$120", Image: "https://unsplash.com/photos/urjasxHT9Ck/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8OHx8bHVtYmVyfGVufDB8fHx8MTc1NjU1NjcwNnwy&force=true&w=128", Location: "Granville, NSW", TimeAgo: "2d ago"},
		},
		Impact: []models.ImpactItem{
			{Label: "Contributions", Value: "47"},
			{Label: "Discussions", Value: "23"},
			{Label: "Circle Tithe", Value: "$5/month"},
		},
	}
}

func GetMockProfileData(handle string, isOwner bool) models.ProfileData {
	return models.ProfileData{
		BaseData: models.BaseData{
			Title:     fmt.Sprintf("%s - Profile", handle),
			ActiveNav: "profile",
			Theme: models.ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: utils.GenerateCSRFToken(),
		},
		Profile: models.Profile{
			ID:     "maia",
			Handle: "@maia",
			Name:   "Maia Makes",
			Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=128&h=128&fit=crop&crop=face",
			Banner: "https://images.unsplash.com/photo-1597960194599-22929afc25b1?w=1200&h=300&fit=crop&crop=center",
			Bio:    "Woodworker & furniture maker crafting heirloom pieces from sustainably sourced timber. Teaching traditional joinery techniques and sharing the journey from tree to table.",
			Stats: models.ProfileStats{
				Posts:       3,
				Connections: 342,
				Circles:     5,
			},
			IsConnected: false,
		},
		Posts: []models.Post{
			{
				ID: "1",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2h ago",
				Circle:  "Woodworking",
				Image: &models.MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=600&h=400&fit=crop&crop=center",
					Alt: "Oak coffee table project",
				},
				CanBuy: true,
			},
			{
				ID: "2",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Spending today selecting timber for the next commission. There's something meditative about running your hands along the grain, feeling for the perfect piece that wants to become a dining table. The wood tells its own story - weather marks, growth patterns, all the years it spent reaching toward light.",
				TimeAgo: "1d ago",
				Circle:  "Woodworking",
			},
			{
				ID: "2",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Traditional style - the backbone of solid furniture. Here's the technique I learned from my mentor, passed down through generations of craftspeople. No shortcuts, just sharp tools and patient hands.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
				Video: &models.MediaItem{
					URL: "https://www.pexels.com/download/video/5972633/",
					Alt: "Video showing old school planing of an uneven timber edge",
				},
				Stats: &models.PostStats{
					Replies: 34,
					Shares:  67,
				},
			},
			{
				ID: "3",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "With over a decade in the business, it's nice to still have all my fingers.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
			},
		},
		PostOffset:   1,
		HasMorePosts: true,
		IsOwner:      isOwner,
	}
}

func GetMockProfileInternalData() models.ProfileData {
	return models.ProfileData{
		BaseData: models.BaseData{
			Title:     "My Profile",
			ActiveNav: "profile",
			Theme: models.ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: utils.GenerateCSRFToken(),
		},
		Profile: models.Profile{
			ID:     "maia",
			Handle: "@maia",
			Name:   "Maia Makes",
			Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=128&h=128&fit=crop&crop=face",
			Banner: "https://images.unsplash.com/photo-1597960194599-22929afc25b1?w=1200&h=300&fit=crop&crop=center",
			Bio:    "Woodworker & furniture maker crafting heirloom pieces from sustainably sourced timber. Teaching traditional joinery techniques and sharing the journey from tree to table.",
			Stats: models.ProfileStats{
				Posts:       3,
				Connections: 342,
				Circles:     5,
			},
			IsConnected: false,
		},
		Extensions: []models.Extension{
			{ID: "pos", Name: "Point of Sale", Description: "Sell directly from your profile", Enabled: true},
			{ID: "analytics", Name: "Analytics", Description: "Get insights about your posts", Enabled: true},
			{ID: "polls", Name: "Polls & Surveys", Description: "Create community polls", Enabled: false},
			{ID: "booking", Name: "Schedules and Booking", Description: "Manage a schedule for booking a service", Enabled: true},
			{ID: "skills", Name: "Skill Exchange", Description: "Offer and request skills", Enabled: false},
			{ID: "gallery", Name: "Portfolio Gallery", Description: "A customisable and currated showcase", Enabled: true},
		},
		Analytics: models.Analytics{
			ProfileViews:         1247,
			ProfileViewsChange:   12,
			PostEngagement:       89,
			PostEngagementChange: 5,
			NewConnections:       23,
			NewConnectionsChange: -3,
		},
		Drafts: []models.DraftPost{
			{
				ID:      "draft1",
				Content: "Workshop tour for anyone thinking of booking some time to create their next masterpiece. Also sharing some essential tools and layout considerations I've learned over 15 years of making.",
				Gallery: []models.MediaItem{
					{URL: "https://unsplash.com/photos/0CCVIuAjORE/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8N3x8d29ya3Nob3AlMjB3b29kfGVufDB8fHx8MTc1NjY5MzI3M3ww&force=true&w=640", Alt: "Workshop space"},
					{URL: "https://unsplash.com/photos/cSqDUEBQUAQ/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8Nnx8d29ya3Nob3AlMjB3b29kfGVufDB8fHx8MTc1NjY5MzI3M3ww&force=true&w=150&h=120&fit=crop&crop=center", Alt: "Tool rack"},
					{URL: "https://unsplash.com/photos/PC9EDk5aDtc/download?ixid=M3wxMjA3fDB8MXxzZWFyY2h8MTR8fHdvcmtzaG9wJTIwd29vZHxlbnwwfHx8fDE3NTY2OTMyNzN8MA&force=true&w=150&h=120&fit=crop&crop=center", Alt: "Wood storage"},
				},
			},
		},
		DraftCount: 1,
		Posts: []models.Post{
			{
				ID: "1",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Just finished this oak coffee table! Happy to step out of my comfort-zone and share some joinery! This piece is available 💜💸",
				TimeAgo: "2h ago",
				Circle:  "Woodworking",
				Image: &models.MediaItem{
					URL: "https://images.unsplash.com/photo-1707749522150-e3b1b5f3e079?w=600&h=400&fit=crop&crop=center",
					Alt: "Oak coffee table project",
				},
				Stats: &models.PostStats{
					Replies: 23,
					Shares:  47,
					Views:   156,
				},
				CanBuy: true,
			},
			{
				ID: "2",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Spending today selecting timber for the next commission. There's something meditative about running your hands along the grain, feeling for the perfect piece that wants to become a dining table. The wood tells its own story - weather marks, growth patterns, the years it spent reaching toward light.",
				TimeAgo: "1d ago",
				Circle:  "Woodworking",
				Stats: &models.PostStats{
					Replies: 12,
					Shares:  28,
					Views:   89,
				},
			},
			{
				ID: "3",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "Traditional style - the backbone of solid furniture. Here's the technique I learned from my mentor, passed down through generations of craftspeople. No shortcuts, just sharp tools and patient hands.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
				Video: &models.MediaItem{
					URL: "https://www.pexels.com/download/video/5972633/",
					Alt: "Video showing old school planing of an uneven timber edge",
				},
				Stats: &models.PostStats{
					Replies: 34,
					Shares:  67,
				},
			},
			{
				ID: "4",
				User: models.User{
					ID:     "maia",
					Handle: "@maia",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=32&h=32&fit=crop&crop=face",
				},
				Content: "With over a decade in the business, it's nice to still have all my fingers.",
				TimeAgo: "3d ago",
				Circle:  "Woodworking",
				Stats: &models.PostStats{
					Replies: 16,
					Shares:  5,
				},
			},
		},
		PostOffset:   3,
		HasMorePosts: false,
		IsOwner:      true,
	}
}
