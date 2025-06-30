# 🐾 The Prowler: Meowl's Web Spider

Every good search engine starts with a curious cat, and for Meowl, that's our web spider! This component is the adventurous pioneer, venturing out into the vast, wild internet to sniff out and gather all the interesting information we need to build our majestic index. It's built to explore systematically, just like a cat exploring every nook and cranny of a new house.

## What's This Kitty Doing?

The Meowl Spider's primary mission is to efficiently fetch web pages and extract the juicy bits of information from them.

- **Breadth-First Search (BFS) Traversal:** Our spider explores the web layer by layer. It prioritizes visiting all known links at a certain "depth" before diving deeper. This ensures a broad coverage of the web before getting lost in the endless tunnels of a single domain.
- **URL Deduplication:** Nobody likes seeing the same mouse twice! The spider keeps track of all URLs it has seen (and is about to see) to avoid redundant crawls and ensure efficiency.
- **Polite Prowling:** We're friendly felines! The spider is designed to respect `robots.txt` rules and implement crawl delays to avoid overwhelming websites.
- **HTML Parsing & Text Extraction:** Once a page is fetched, the spider carefully parses the HTML to extract the visible text content, along with important metadata like titles, descriptions, and links.
- **Link Discovery:** It meticulously uncovers all internal and external links found on the page, adding new potential adventures to its to-do list (the crawl frontier).
- **Data Storage:** All the precious treasures (URL, title, description, raw content, links, etc.) are neatly cataloged and stored in our database, ready for the indexer's hungry paws.
- **Concurrent Crawling (Goroutines & Mutex):** To be as agile as a cat, the spider leverages Go's goroutines for concurrent fetching, making sure it's always busy without stepping on its own tail (managed by mutexes for safe shared resource access).

## How it Hunts

The spider starts with a few initial "seed" URLs and then expands its hunt from there:

1.  **Seeds of Adventure:** It begins with a list of initial URLs to visit (our favorite scratching posts!).
2.  **The Frontier:** These initial URLs, and all subsequent ones discovered, are added to a "crawl frontier" – a queue of URLs waiting to be explored.
3.  **Fetch & Parse:** One by one (or concurrently, thanks to goroutines!), the spider fetches a URL from the frontier, reads its `robots.txt`, downloads the HTML, and parses out the relevant content and links.
4.  **Dedup & Queue:** Newly discovered links are checked against the list of already-seen URLs. If new, they are added to the frontier.
5.  **Store the Catch:** The extracted data (content, title, description, links, etc.) for each page is then saved to our database, complete with a timestamp of when it was caught.
6.  **Repeat:** The process continues until the desired number of links are crawled, a specific depth is reached, or no new links are found.
