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
			{ID: "3", Title: "Oak Lumber Bundle", Price: "$120", Image: "https://images.unsplash.com/photo-1702195789139-4897ff9b0083&force=true&w=128", Location: "Granville, NSW", TimeAgo: "2d ago"},
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

func GetMockCirclesPageData() models.CirclesPageData {
	return models.CirclesPageData{
		BaseData: models.BaseData{
			Title:     "My Circles",
			ActiveNav: "circles",
			Theme: models.ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: utils.GenerateCSRFToken(),
		},
		Circles: []models.Circle{
			{
				ID:           "1",
				Name:         "Woodworking",
				Description:  "Traditional craftsmanship meets modern techniques. Share projects, ask questions, and connect with fellow makers.",
				Thumbnail:    "https://images.unsplash.com/photo-1702195789139-4897ff9b0083?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1497219055242-93359eeed651?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "47",
				OnlineCount:  "8",
				UserRole:     "admin",
				JoinedDate:   "6 months ago",
				LastActivity: "2m ago",
				Active:       true,
			},
			{
				ID:           "2",
				Name:         "The Crop Circle",
				Description:  "The official circle for Alienated Collective, a sydney based group of local ravers, doofers, party and music enthusiasts.",
				Thumbnail:    "https://images.unsplash.com/photo-1611402858501-d3de70f3c67e?q=80&w=80&h=80&auto=format&fit=crop",
				Banner:       "https://images.unsplash.com/photo-1594623930572-300a3011d9ae?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "12",
				OnlineCount:  "3",
				UserRole:     "member",
				JoinedDate:   "3 months ago",
				LastActivity: "15m ago",
				Active:       true,
			},
			{
				ID:           "3",
				Name:         "DIY Electronics",
				Description:  "Arduino, Raspberry Pi, circuit design, and everything in between. From beginner tutorials to advanced projects.",
				Thumbnail:    "https://images.unsplash.com/photo-1603732551658-5fabbafa84eb?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "156",
				OnlineCount:  "24",
				UserRole:     "member",
				JoinedDate:   "8 months ago",
				LastActivity: "1h ago",
				Active:       true,
			},
			{
				ID:           "4",
				Name:         "Sydney Artists",
				Description:  "Supporting Sydney's creative community through collaboration, critique, and celebration of local talent.",
				Thumbnail:    "https://images.unsplash.com/photo-1541961017774-22349e4a1262?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1663505819040-00bbd0814fab?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "89",
				OnlineCount:  "12",
				UserRole:     "member",
				JoinedDate:   "4 months ago",
				LastActivity: "2h ago",
				Active:       false,
			},
			{
				ID:           "5",
				Name:         "Communication Software",
				Description:  "Discussing the future of digital communication, platform design, and better community building tools.",
				Thumbnail:    "https://images.unsplash.com/photo-1526045612212-70caf35c14df?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "34",
				OnlineCount:  "3",
				UserRole:     "owner",
				JoinedDate:   "1 year ago",
				LastActivity: "1h ago",
				Active:       true,
			},
		},
		RecentActivity: []models.CircleActivity{
			{
				ID:       "1",
				CircleID: "1",
				Type:     "post",
				Title:    "New oak coffee table project completed",
				Content:  "Just finished this oak coffee table! Happy to step out of my comfort-zone...",
				User:     "@maia",
				TimeAgo:  "2m ago",
			},
			{
				ID:       "2",
				CircleID: "2",
				Type:     "event",
				Title:    "BBQ gathering tonight",
				Content:  "Warming up the barbeque and just got couple cases of the finest bread-water...",
				User:     "@heathtyler",
				TimeAgo:  "15m ago",
			},
			{
				ID:       "3",
				CircleID: "3",
				Type:     "announcement",
				Title:    "New Arduino tutorial series starting",
				Content:  "First session covers basic circuits and programming fundamentals...",
				User:     "@sara_pcb",
				TimeAgo:  "1h ago",
			},
			{
				ID:       "4",
				CircleID: "4",
				Type:     "member_joined",
				Title:    "New member joined Sydney Artists",
				Content:  "@alex_painter has joined the circle",
				User:     "System",
				TimeAgo:  "3h ago",
			},
			{
				ID:       "5",
				CircleID: "1",
				Type:     "post",
				Title:    "Workshop layout considerations",
				Content:  "Spending today selecting timber for the next commission...",
				User:     "@maia",
				TimeAgo:  "1d ago",
			},
		},
		Stats: models.CircleStats{
			TotalPosts:     127,
			ActiveMembers:  52,
			RecentActivity: "Very High",
			WeeklyGrowth:   "+12%",
			EngagementRate: "85%",
		},
		FeaturedCircles: []models.Circle{
			{
				ID:           "6",
				Name:         "Sustainable Living",
				Description:  "Practical tips for reducing environmental impact through daily choices and community action.",
				Thumbnail:    "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "234",
				OnlineCount:  "18",
				UserRole:     "",
				JoinedDate:   "",
				LastActivity: "",
				Active:       true,
			},
			{
				ID:           "7",
				Name:         "Urban Beekeeping",
				Description:  "City-based beekeeping community sharing techniques, equipment, and harvest stories.",
				Thumbnail:    "https://images.unsplash.com/photo-1558642452-9d2a7deb7f62?w=80&h=80&fit=crop&crop=center",
				Banner:       "https://images.unsplash.com/photo-1558642452-9d2a7deb7f62?w=400&h=200&fit=crop&crop=center",
				MemberCount:  "78",
				OnlineCount:  "6",
				UserRole:     "",
				JoinedDate:   "",
				LastActivity: "",
				Active:       true,
			},
		},
	}
}

func GetMockChatData() models.ChatPageData {
	return models.ChatPageData{
		BaseData: models.BaseData{
			Title:     "Chat",
			ActiveNav: "chat",
			Theme: models.ThemeSettings{
				Mode:   "system",
				Radius: "0",
			},
			CSRFToken: utils.GenerateCSRFToken(),
		},
		Conversations: []models.Conversation{
			{
				ID:          "1",
				Name:        "circles.diy Design",
				Avatar:      "https://images.unsplash.com/photo-1581291518857-4e27b48ff24e?w=48&h=48&fit=crop",
				LastMessage: "Maya: The color palette looks great! 🎨",
				LastTime:    "2m ago",
				UnreadCount: 3,
				IsOnline:    true,
				IsGroup:     true,
				Participants: []models.User{
					{ID: "maya", Handle: "@maya", Name: "Maya Chen", Avatar: "https://images.unsplash.com/photo-1522075469751-3847ae96c725?w=32&h=32&fit=crop&crop=face"},
					{ID: "alex", Handle: "@alex", Name: "Alex Ramirez", Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face"},
					{ID: "jordan", Handle: "@jordan", Name: "Jordan Kim", Avatar: "https://images.unsplash.com/photo-1494790108755-2616b2e6ead5?w=32&h=32&fit=crop&crop=face"},
				},
			},
			{
				ID:          "2",
				Name:        "Emma Wilson",
				Avatar:      "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=48&h=48&fit=crop&crop=face",
				LastMessage: "Thanks for the feedback on the prototype!",
				LastTime:    "15m ago",
				UnreadCount: 0,
				IsOnline:    true,
				IsGroup:     false,
			},
			{
				ID:          "3",
				Name:        "Dev Team",
				Avatar:      "https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=48&h=48&fit=crop&crop=face",
				LastMessage: "Sam: Ready for the demo tomorrow? 💻",
				LastTime:    "1h ago",
				UnreadCount: 1,
				IsOnline:    false,
				IsGroup:     true,
				Participants: []models.User{
					{ID: "sam", Handle: "@sam", Name: "Sam Rodriguez", Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face"},
					{ID: "riley", Handle: "@riley", Name: "Riley Park", Avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=32&h=32&fit=crop&crop=face"},
				},
			},
			{
				ID:          "4",
				Name:        "Marcus Thompson",
				Avatar:      "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=48&h=48&fit=crop&crop=face",
				LastMessage: "Great job on the presentation!",
				LastTime:    "2h ago",
				UnreadCount: 0,
				IsOnline:    false,
				IsGroup:     false,
			},
			{
				ID:          "5",
				Name:        "Book Club Circle",
				Avatar:      "https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=48&h=48&fit=crop&crop=face",
				LastMessage: "Lisa: Anyone else loving this chapter? 📖",
				LastTime:    "3h ago",
				UnreadCount: 0,
				IsOnline:    true,
				IsGroup:     true,
			},
		},
		ActiveChat: &models.Conversation{
			ID:       "1",
			Name:     "circles.diy Design",
			Avatar:   "https://images.unsplash.com/photo-1581291518857-4e27b48ff24e?w=48&h=48&fit=crop&crop=face",
			IsOnline: true,
			IsGroup:  true,
			Participants: []models.User{
				{ID: "maya", Handle: "@maya", Name: "Maya Chen", Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=48&h=48&fit=crop&crop=face"},
				{ID: "alex", Handle: "@alex", Name: "Alex Ramirez", Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face"},
				{ID: "jordan", Handle: "@jordan", Name: "Jordan Kim", Avatar: "https://images.unsplash.com/photo-1494790108755-2616b2e6ead5?w=32&h=32&fit=crop&crop=face"},
			},
		},
		Messages: []models.Message{
			{
				ID:        "1",
				Content:   "Hey everyone! I've been working on the new design system. What do you think about these color combinations?",
				Timestamp: "10:30 AM",
				Sender: models.User{
					ID:     "maya",
					Handle: "@maya",
					Name:   "Maya Chen",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=48&h=48&fit=crop&crop=face",
				},
				IsOwn:  false,
				IsRead: true,
				Type:   "text",
			},
			{
				ID:        "2",
				Content:   "Those shades are perfect! Really captures the minimal vibe.",
				Timestamp: "10:32 AM",
				Sender: models.User{
					ID:     "alex",
					Handle: "@alex",
					Name:   "Alex Ramirez",
					Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face",
				},
				IsOwn:  false,
				IsRead: true,
				Type:   "text",
			},
			{
				ID:        "3",
				Content:   "I love the direction this is taking! The contrast ratios look accessible af 🔥",
				Timestamp: "10:35 AM",
				Sender: models.User{
					ID:     "current_user",
					Handle: "@you",
					Name:   "You",
					Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face",
				},
				IsOwn:  true,
				IsRead: true,
				Type:   "text",
			},
			{
				ID:        "5",
				Content:   "Perfect! This is exactly what I had in mind. Should we schedule a call to discuss implementation?",
				Timestamp: "10:45 AM",
				Sender: models.User{
					ID:     "jordan",
					Handle: "@jordan",
					Name:   "Jordan Kim",
					Avatar: "https://images.unsplash.com/photo-1543610892-0b1f7e6d8ac1?w=32&h=32&fit=crop&crop=face",
				},
				IsOwn:  false,
				IsRead: true,
				Type:   "text",
			},
			{
				ID:        "6",
				Content:   "Great idea! I'm free this afternoon. How about we gather @ 2 PM?",
				Timestamp: "10:46 AM",
				Sender: models.User{
					ID:     "current_user",
					Handle: "@you",
					Name:   "You",
					Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face",
				},
				IsOwn:  true,
				IsRead: true,
				Type:   "text",
			},
			{
				ID:        "7",
				Content:   "Thanks guys 🫶 Sounds good, talk soon!",
				Timestamp: "10:48 AM",
				Sender: models.User{
					ID:     "maya",
					Handle: "@maya",
					Name:   "Maya Chen",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=48&h=48&fit=crop&crop=face",
				},
				IsOwn:  false,
				IsRead: false,
				Type:   "text",
			},
		},
		Contacts: []models.Contact{
			{
				User: models.User{
					ID:     "maya",
					Handle: "@maya",
					Name:   "Maya Chen",
					Avatar: "https://images.unsplash.com/photo-1653508242641-09fdb7339942?w=48&h=48&fit=crop&crop=face",
				},
				IsOnline:     true,
				Relationship: "circle_member",
			},
			{
				User: models.User{
					ID:     "alex",
					Handle: "@alex",
					Name:   "Alex Ramirez",
					Avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=32&h=32&fit=crop&crop=face",
				},
				IsOnline:     true,
				Relationship: "friend",
			},
			{
				User: models.User{
					ID:     "emma",
					Handle: "@emma",
					Name:   "Emma Wilson",
					Avatar: "https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=32&h=32&fit=crop&crop=face",
				},
				IsOnline:     true,
				Relationship: "friend",
			},
			{
				User: models.User{
					ID:     "marcus",
					Handle: "@marcus",
					Name:   "Marcus Thompson",
					Avatar: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=32&h=32&fit=crop&crop=face",
				},
				IsOnline:     false,
				LastSeen:     "2h ago",
				Relationship: "circle_member",
			},
		},
	}
}
