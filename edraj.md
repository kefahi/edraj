# edraj

### The next-generation Information Management System
<img src="revealjs/images/edraj.png" alt="Websocket" width="50%" style="float: left; background:none; border:none; box-shadow:none;">

<p style="text-align: right; width: 100%;">Version 0.3<br/>Amman, January - 2018<br/>kefah.issa@gmail.com</p>

---
# The challenge

<img src="revealjs/images/quote2.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
###	_One app/service to handle most of the user content management and communication needs (private/shared/public) as opposed to the scattered Social media outlets and user managed content options in addition to figuring the key qualities for that Next-gen app to succeed._


+++
### Pain-points
- Users don't really own their content:
	- The older content is quickly buried and hard to re-surface and search.
  - Users can lose content, pages, account when banned / reported hence losing:
    - The content they contributed
    - The community reactions
    - And above all their social graph (network). e.g. Pages with millions of followers can be banned.
- Social Media / Content Management providers are proprietary and centrally owned/managed, they mainly strive to lock-in the users:
  - No API that enable the users to fully manage their content, let alone that federation options are not welcomed by the providers.
  - Users can not automatically and fully **sync** content across various outlets
- Users are overwhelmed
  - With the ever increasing debt of their own content that they may fail to manage and hence lose.
  - With the sheer number of outlets (providers) that overlap and confuse the users

+++
### Current Status-quo  <img src="revealjs/images/openlock.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

A wide range of use-cases served by independent and greatly over-lapping applications each attempting to lock-in the user's data.
- **Social Media apps**: Facebook, Instagram, Twitter
- **Public CMS, Blogs, News Syndication** : Word-press, RSS/Atom
- **Selective-sharable Content/Media store**: e.g. Cloud drives
- **Messaging and Email**
- **_Future_**: Knowledge base (e.g. wiki), note-taking, todos and ticketing
- **_Future_**: Sheets, Inventory: e.g. Excel, Access, Foxpro

+++
### Aspiration <img src="revealjs/images/aspiration.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

One unified set of open-standards that enable all those use-cases (is able to handle various types of content) and run in a federated fashion. Offering superior content-management experience and eliminating vendor-lock-in that is overwhelming the users.

Email (SMPT/IMAP4/POP3) comes as an example of open-standard used by the "Email" use-case.

+++
### Barriers  <img src="revealjs/images/barrier.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* Technically speaking there is no reason why one application (with a proper ecosystem) cannot serve all those use-cases and more. At the end of the day it boils down to information/data entities that should be maintained, indexed and shared / published.
* Unfortunately and despite the fact that there have been several attempts to address sub-sets of use-cases, none of them tries to address the problem in its fully generality.
Let alone offering it at competitive quality and speed.
* The social affinity to the existing social network offering (existing social graph weight).

---
### Crossing the chasm <img src="revealjs/images/cross.svg" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

Significant advantages must be offered for the users to switch:

1. **Unified data and meta-data**: All data/content is well-defined and managed using the same verbs. (e.g. REST, schema.org/JSON-LD)
2. **Federated/self-hosted and Standard-API-based**: Inter-operable independently-hosted domains (DHT) (sets of users). Also key-based idM (e.g. OpenId connect).
	* True account/content ownership
3. **Web-enabled/SEO-friendly availability of content**: Just think how quickly the user-contributed content gets buried and almost unreachable to others.
4. **Smart**: Programmable time/trigger-based activities. e.g. IFTTT
5. **Free/Open-source**: Grant users the right to use, modify and evolve the implementations.


+++
### Crossing the chasm (continued) <img src="revealjs/images/cross.svg" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

6- **Performance and Quality**: The app **must** be very performing and not less stable than the current offerings.

7- **Privacy and Security**

8- **Easy service and data redundancy options**

+++
### Unified data and meta-data <img src="revealjs/images/metaData1.svg" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* File-based content and meta-data persistence. The content and its meta-data file live in the file-system under certain relative path.
* Content is well-defined and self-described (through meta-data)
* The Meta-data for each content include:
  * Name
  * Content _guid_
  * Owner _guid_
  * Access-control settings (permissions): for Actors, Groups and Work-groups.
  * Time-stamps: Created at, Last-updated
  * Change-history for revision-enabled content
  * Labels: Tags, Categories
  * Path? (redundant from the file-system)
  * Schema reference definition (for structured data)

+++
### Content Management Basics  <img src="revealjs/images/metaData.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* Every entity in the system has a life-long globally unique identifier (guid): Actor, Work-group, Content, ...etc.
* Basic **Verbs**:
  * Create
  * Update content / meta-data (Access-control/permission included)
  * Delete
  * Query (with filtering)
  * Subscribe/Poll for notifications
* Every thing is persisted on the file-system including messages. _implementation recommendation_.
* The caching/index component should always be able to completely rebuild the cache/indexes as such it should only be considered for performance purposes.
* Future: both in-motion and in-transit data should be encrypted

+++
### Federated and Standard-API-based <img src="revealjs/images/federated.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

- A Domain is a collection of users along with the various edraj components (see <a href="#/5">Components section</a>).
- Federation enables users and groups / communities to self-host and inter-operate with everyone else.
- Federation and Standard-API are key to free users from vendor-lock-in. It even allows multiple-technical implementations.
- Federation is also a means to eliminate the concept of one single service-provider, helping distribute the processing load to a manageable level as opposed to requiring huge investment in infrastructure. As it eliminates the need / requirement for one central service provider; it will simply distribute the cost of hosting (computing/storage/data-transfer) over the federated domains.

+++
### Federated and Standard-API-based (continued) <img src="revealjs/images/federated.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

Federation and the option of self-hosted inter-operable information management system empowered the users and communities to **finally** and **truly** own their content.
  - No Tracking
  - No Advertisement
  - No one can delete their files (posts/videos ...)
  - No one can suspend their account
  - The user no longer needs to beg YouTube, Facebook and twitter not to action against them.
  - Yes, there is a financial element associated with self-hosting; but its a small-price to pay that would set you free.

+++
### Performance and Quality <img src="revealjs/images/quality.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

Imagine how poorly would an app be received if it lacks either of those two qualities. No matter how great its idea is, how nice-looking, how easy, ...etc. As such those two qualities are very key and are actually make-or-break.
---

### Key concepts and definitions  <img src="revealjs/images/definitions.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* **Actor**: Individual agent / user.
* **Work-group**: Collaboration of a group of actors on workgroup-owned content (like FB pages/groups)
* **Content**: Media, Document, Text + self-describing Meta data.
	* **Labels**: Tags (user-defined, free-form) and Categories (predefined, hierarchical)
	* **Container**: Hierarchical aggregation of content: folder/tar-ball
	* **Permission**: Privileges granted to owner and people
* **Action**: Manage content, React to content, subscribe to notifications (filter-based actor/tag/category/type),
* **Notification**: Event notifications for subscribed users.
	* **Message**: One-to-one or one-to-many messaging (think Email and Instant messaging)
* **Page**: A single page view made of **Blocks** (physical view-able manifestation of the content).

+++
### Actor <img src="revealjs/images/user.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

User/profile (single-user):
- Information about the user that includes:
  - Bio
  - Key pair(s)
  - Email(s)
  - Phone (contact)
  - Groups of Actors (like access-groups in Facebook) to help simply access-control management
  - Guids/Public-keys/API-url-pointers of other associated users (contacts, followed, friends, ...etc)

+++
### Work-group  <img src="revealjs/images/workgroup.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

- Team-centeric (communities/interest groups) collaborative content management and communication medium.
- Members from any domain can join/be invited per the privacy settings.

+++
### Content <img src="revealjs/images/contentMang.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
* Owners can: 
	- Admin (Meta/privileges)
	- Create/update/delete
	- Manage comments
* Labels are used to qualify the content for better search and organization.
* Content is also organized in arbitrary folder-structures, like regular file-systems and cloud-drives.
* Future: Large content can be served in a torrent-like fashion.


+++
### Action <img src="revealjs/images/action.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* Manage content: Create/Update/Delete, Set meta-data options (including permissions)
* React: comment, share, like/dislike/interesting ...etc.
* Message / Communicate
* Subscribe to notification based on a filter.

+++
### Notification <img src="revealjs/images/notifications3.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

Notification-streams (follow/unfollow person or work-group or specific-content, tag)
A user can follow notifications (get stream of notifications) based on filter-rules:

- Actor Activities
- Work-group Activities
- Tags / Categories (Horizontally trending general public content) twitter-hashtag like (What's trending)
- Specific types of notifications
- The ranking formula is also affected by positive reactions (likes + rating), shares, comments and people who reacted.

+++
### Page / Block <img src="revealjs/images/pages.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

A page is a public presentation of selected content / streams. (blog-sphere like)

* Users can choose layout/template and how content is presented by setting queries. with pagination support.
* The layout groups a number of blocks
* Each block has a query-filter to determine the content to be surfaced and a template to determine how it is presented.


---
<img src="revealjs/images/quote2.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
##	_For every respective set of use-cases  edraj can be reduced to an existing known / popular platform: Blog, Email, Messaging, Social media, News services, Media management._

---
### Main components <img src="revealjs/images/puzzle.svg" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* **Identity Manager (idM)**: Manage Actors/users and people groups.
* **Content Manager**: Create/Manage, React, and publish (permission / access control)
* **Notifications and Messaging**: Action-driven notification system and messaging capability.
	* **_Future_: Peer-to-peer**: Audio/Video conferencing
* **Miner**: Index and attempt to improve meta-data. also remote-index (gopher-like)
* **Add-ons**: local and remote repo of usable add-ons (mini-apps): e.g. Interface with existing social media sites: import/export (sync).
* **Public interface**: Public content including permissions, pages and blocks + notification subscription/distribution (syndication).
* **Schema definitions**: structured data: standard + custom: Local + remote repos
* **Client-apps**: Mobile, Desktop, Web-SPA

+++
### Identity Manager (idM)  <img src="revealjs/images/id.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* Setup per domain, hosting an arbitrary number of users as determined by the admin
* The identity manager also stores the public keys (and links to them) of users on other domains (followed / friends).
* Think OpenID connect (OAuth2) compliant.


+++
### Content Manager <img src="revealjs/images/content.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

- Persistence of data / content with the proper Meta-data
- Types of content:
  - Plain text
  - Rich text
  - Structured data (person json schema)
  - Message
  - Wiki
  - Binary-payload: Payload of Media files, Documents or any other type of binary attachment (supports both embedded and URL-pointers ).
- Types of container:
  - Folder (regular file-system folder)
  - Tar-ball (compressed hierarchical folder/file structure).


+++
### Miner  <img src="revealjs/images/mining.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

- Index and classification of local data that enriches local meta data
- Indexing of public content from other domains (public and permitted content)
- Polling other domains for notifications (a client may only be notified when a domain has something new, then it should poll for it).
- Public mining-only services could exist that would horizontally aggregate and promote relevant trending content.

+++
### Add-ons <img src="revealjs/images/addon.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
- Plugins/Mini-apps
- Public repositories (one formal)

+++
### Public interface <img src="revealjs/images/interface.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

- Serves public web-enabled/SEO-friendly content: Pages + blocks, including generating proper robots.txt+sitemaps.
- Serves public API for
	* Notification-polling and subscription
  * Content Query and retrieval
  * Surface Content management component public api's
  * Messaging and communication
  * Structured content query (schema-enabled content in a machine-consumable form)

+++
### Schema definitions <img src="revealjs/images/database.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
Predefined schema definitions allow machine consumption of properly described structured data.
The local + remote repos that contain the schema definitions for structured data.

Examples of data that require structured schema definitions:
- Task/Todo-entry
- Contact details: Person or Organization
- Scheduled Event: Public/private event, meeting, reminder ...
- Place + tracking
- Blog/Post/Short-post(tweet-like)
- Term / phrase definition or translation
- Quote / Proverb
- Biography

Each one of these data types is best represented by its own standard schema definition.

+++
### Client-apps <img src="revealjs/images/app.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
- Mobile apps: Progressive Android/iOS ..etc.
- Web/Desktop apps: Progressive Web/Desktop apps and/or Native apps.

+++
### Distributed hashtables (DHT) <img src="revealjs/images/consul.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">
One official DHT exists and is configured per default, additional ones can be configured.
- Domains: which are not DNS-based for two reasons:
	- Enable more liberal naming (names in the formal repo are reviewed)
	- Reduce the setup complexity which opens doors for behind-firewall peer-to-peer intr-domain communication.
- Add-ons
- Schema
- Public miners (search-engine like)

---
### High-level Architecture diagram 
<img src="revealjs/images/edraj-deployment.svg" style="background:none; border:none; box-shadow:none;" width="100%" >

+++
### Standards definition (OpenApi 3) _Work-in-progress_ <img src="revealjs/images/openapi-icon.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

<a href="swagger-ui-dist/index.html" target="_blank">click here for edraj api specifications</a>

+++
### Storage abstraction notes

* json/meta-data enabled, path-aware file-store : basic backend is filesystem-based.
* Multiple isolated roots: domains, addons, schema, messages, notifications, content, local-miner, public-miners, public-serving: pages/blocks/layouts/templates/js/images/css/links to public content-payload, trash, content-cache (for all externally accessed content), curated: external content that the user reacted to.
* file-name, .file-name.meta.json (can be stored in XFS file-attributes?), folder:.file-name.changes
* Future: Reference checking capability? 
	- when a schema definition is removed, no file-meta should be pointing to it.
	- When a file is attached to a message?
* Future: Main-feature: automatic deduping: used hart/soft? links

+++
### Main storage roots

* **content**:
	*	Media files (video, audio, images)
	* Documents / ebooks (readable material)
	* Text (plain / rich)
  * Links (urls bookmarks / reddit like)
	* Structured data
* **people**: subscribed-to people (followed/friends) + their groupings
* **messages**: Message store (instant and email)
* **notifications**: Events on interesting activities
* **public-serving**: Pages, Blocks, Layouts, Templates, Static files: JS/CSS/Images
* **local-miner**: Index area of the local and cached data. Rebuildable from scratch.
* **domains**: cached domain-details
* **schema**: schema used by other parts of the local setup
* **addons**: installed modules/addons
* **trash**: Where deleted data is moved until permenantly deleted.


+++
### Main Folder structure
	├── content
	├── messages
	├── notifications (events)
	├── public-facing (serving public pages)
	├── people
	├── indexes (local-miner)
	├── modules (addons)
	├── domains
	├── schema
	└── trash

+++
### Suggested "content" structure

This is user-managed

	content
		├── personal (private)
		├── family
		├── friends
		├── interests
		│	└── links
		├── messages
		├── files
		│	├── books and documents
		│	├── media
		│	└── apps
		└── structured-data
			├── tickets (and todos)
			├── inventory
			├── study
			└── financial

+++
### Storage abstraction api

1. **PutMeta**:  Full or Delta (only mentioned fields are updated). Data integrity is checked before and after the put-operations. creates the folder-path if it doesn't exist.
2. **GetMeta**: json formatted meta (full or sub-set)
3. **List**: folder-only, Returns 1st-level child file and folder meta data.
4. **PutPayload**: file-only. Full / patch?(binary/text) + checksum
5. **GetPayload**: file-only
6. **Delete**: recurse-option (i.e. delete belongings/childern -> deletes both meta-data and otherwise; must be reversable (use some trash concept))
7. **Move**: Aka rename. within the same parent folder or to another one.

+++
### Misc. notes

* When content is shared (video played) (or interacted-with) its also copied locally and served in a torrent-like fashion. (i.e. each interaction is registered as a service point) a DHT is needed to serve this: Curated content DHT. ( an asynchronus connection with edraj should always be openned. /Websocket/QUIC/Http2(gRPC) based? (two-way-stream)
* A content has multiple-views: isolated json/binary, isolated embeddable, isolated web, future: parent web?

---

### Roadmap and milestones <img src="revealjs/images/milestone2.png" alt="Imagine yourself" width="20%" style="float: right; background:none; border:none; box-shadow:none;">

* **Phase I. Conception** : December/17
* **Phase II. Prototype** : June/18
* **Phase III. MVP first release**  : December/18
* **Phase IV. Continuous releasing and improvement**
