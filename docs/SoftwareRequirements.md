# Crawler Requirements Document 



## 1.0 Introduction

The Internet Archive in collaboration with several National Libraries is seeking to build an open source Crawler that can be used primarily for Web Archiving purposes meeting the requirements defined by the National Libraries Web Archive Consortium.

This document describes in part the anticipated requirements for such crawling software. The first subsection discusses the use cases we wish to support (that is, it describes, at a high-level, our anticipated patterns of use). The following section provides a statement of requirements derived from those use cases.

## 2 Use-case analysis

The Archive wants to support the following use cases:

### 2.1 Broad crawling

Broad crawls are large, high-bandwidth crawls in which the number of sites and the number of interesting, individual pages touched is as important as the completeness with which any one site is covered. The Web is uncataloged; broad crawling is important for discovering unknown pages. Also, one cannot know _a priori_ what parts of the Web will be important in the future; broad crawling attempts to ensure sufficient coverage for these unknown uses.

### 2.2 Narrow crawling

Narrow crawls are small- to medium-sized crawls in which the quality criterion is complete coverage of some selected sites or topics. Narrow crawls are necessary for capturing sites that are known _a priori_ to be important (either intrinsically or in relation to some collection being put together).

### 2.3 Continuous crawling

Traditionally, crawlers have taken pains to download only pages they have not previously downloaded. Such crawlers are typically known as _batch_ or _snapshot crawlers_ because they are typically run for a fixed amount of time to take a snapshot of the Web, and then are restarted to do a full crawl all over again.

Continuous crawling continuously downloads previously fetched pages (looking for changes) as well as (possibly) discovering and fetching new pages. Continuous crawling is useful for capturing pages that change quickly. Also, when the download rates of pages are adjusted to reflect measured change frequencies, continuous crawling can be more efficient than batch crawling.

The Archive is interested in continuous crawling, both for its broad and its narrow crawls.

### 2.4 Multi-protocol crawling

The Archive is interested in crawling pages available through protocols other than HTTP. RSS is of particular interest. Additional protocols include HTTPS, FTP, and streaming protocols such as MMS, RTSP, and PNA.

### 2.5 Experimental crawling

The Archive wants to support experiments in crawling techniques and technology performed by external organizations. Crawling for the National Digital Science Library done at Cornell University is a good example of what we want to support. In this project, researchers wanted to build a topic-focused crawler to find science-related material on the Web. They were able to write analyzers that used term-vector techniques to identify material of interest and were able to modify the policies of the crawler to use this information to direct the crawl.

We want to be able to support projects such as the NSDL project. To do so, the crawler needs to be easy to extend and easy to use, and it cannot be over-encumbered with licensing constraints.

### 2.6 Zero administration/ease of operation

The Archive does not have a large operations staff and thus requires a high level of automation in its crawlers. In particular, we imagine that a single, broad crawl can be administered by the part-time efforts of a single person (assuming reliable hardware and network connectivity). Also, the Archive cannot provide much technical support to external organizations using its crawling software, so this software must be easy to operate without much support.

## 3 Requirements

We believe that the following requirements are necessary to support the above use cases:

## 3.1 Extensibility

Extensibility is the ability to change the behavior of the crawler: its selection policy, its scheduling policy, what it stores to disk, and so forth. Extensibility is our most important requirement. It is needed to support almost all of the use cases above (broad and narrow crawling, continuous crawling, multi-protocol crawling, ease of operation, and, as mentioned above, experimental crawling).

As enumerated below, extensibility applies to many aspects of the crawler. At a high level, what we really want is a _framework_ for building crawlers, not a particular crawler embodying specific policies or capabilities. A _crawl designer_ should be able to put together an instance of this framework to meet the needs of a specific crawling project. At the same time, we do _not_ want a traditional object-oriented framework that _requires_ programming to put together useful systems. As we discuss below, while the ability to plug new, custom code into the framework is important, it is also important that the framework support extensibility through parameterization as well.

### 3.1.1 Extensible aspects of the crawler

Many different aspects of the crawler need to be extensible. In this subsection, we list _what_ needs to be extensible in the crawler; in the following subsection, we list _how_ this extensibility is to be achieved.

#### 3.1.1.1 Protocols

It should be easy to add new downloading protocols to the framework.

#### 3.1.1.2 Selection (discovery and filtering)

The _selection policy_ determines what the crawler will download. Selection is typically the product of two mechanisms: a _discovery_ mechanism by which new URLs are identified, and a _filtering_ mechanism that decides which of the discovered URLs is to be downloaded. In continuous crawls in which old URLs are recrawled, there also needs to be a filtering process to eliminate URLs that are no longer of interest.

Selection policies are of vital importance to both broad and narrow crawlers. It is an area in which we expect there to be significant, on-going research and experimentation. Thus, our crawler framework must support the widest-possible extensibility when it comes to discovery and filtering mechanisms.

The discovery process can sometimes happen outside of the crawler itself. For example, some organizations have access to Web proxies or other usage information that can be good sources of URLs to crawl. To support external discovery and filtering mechanisms, it must be possible to both add and remove large numbers of URLs from the crawler while a crawl is in process.

#### 3.1.1.3 Scheduling

The _scheduling policy_ (sometimes called _ranking policy_) determines the order in which pending URLs are downloaded. The scheduling and selection policies together define the contents of the collection obtained by a crawl, and thus are of enormous interest to crawl designers. Significant flexibility is required with respect to both. However, flexibility with respect to scheduling is particularly difficult to provide because the scheduling policy is determined by the crawlers _frontier_ of pending URLs.

A crawlers frontier is the collection of URLs that the crawler intends to download in the future. In a batch crawl, the frontier contains only new URLs (URLs that have not been downloaded by this crawl). In a continuous crawl, the frontier can contain both new URLs and old URLs that are due to be recrawled. The scheduling policy determines the order in which elements are removed from this frontier. For broad crawling the frontier can become very large, containing as many as 10s of billions of URLS, and thus, in a scalable crawler, the frontier must be a disk-based data structure. Disk-based data structures can have high latency, but the frontier is on the critical path of the crawler, where latency cannot be tolerated. Thus, a good frontier requires significant programming skill to implement and tune.

As a result, we cannot expect new frontiers to be written quickly and easily, and thus flexibility for scheduling policies cannot come primarily from writing new frontier implementations. Rather, it is important that the crawler framework come with a variety of frontier implementations supporting a variety of scheduling policies. At the same time, to support experimentation and further development, there must be a good API for defining new frontier implementations.

#### 3.1.1.4 Politeness

High-capacity crawlers can easily overwhelm a low-end Web server and can even put stress on high-end Web servers. Thus, it is very important that crawlers be polite, limiting the load that they put onto a Web server.

Unfortunately, the need to be polite sometimes conflicts with the scheduling policy of a crawl. For example, a simple politeness policy is to stagger, over time, the pages downloaded from a given site. However, as discussed above, it is sometimes desirable to download images and other objects associated with a page immediately after the page itself is downloaded. If this latter requirement is particularly important to the crawl, then a more sophisticated politeness policy is needed to accommodate it. At the same time, when such freshness issues are not important to a crawl, it is generally considered neighborly to space downloads from a site as far apart as possible.

The crawler must be flexible regarding the politeness policies it implements. Again, there really is no such thing as a one size fits all crawler; what is needed is a framework that accommodates a wide variety of policies.

#### 3.1.1.5 Document processing

Once a document is downloaded, it is processed for a variety of reasons. Features might be extracted for the purpose of further crawling (e.g., links for most crawls, plus link-text for topic-directed crawls). Features might be extracted for the sake of collecting statistics, and the page itself might be transformed into some appropriate output format.

Thus, the crawler must support a framework for pluggable analyzers that process documents. Further, this framework must support effective interoperability among analyzers, ensuring that the good work of one group can be built-upon by others.

Analyzers should have access to as much metadata as possible. This includes the URL used to fetch the object, download timestamps, DNS resolution information, and request and response headers used during the transaction (plus any other material used to fetch the object, such as passwords or SSL keys).

Analyzers should have access to the results of _all_ requests, not just successful ones. This is necessary, for example, for detecting change in the Web, or for detecting the edges of the deep Web. This means recording, for example, 30x (moved/redirect) and 40x (invalid/not found/unauthorized) responses. It also means recording requests that are prevented by robot exclusion rules. Finally, automatically fetched robot-exclusion files should also be subject to regular analysis (rather than be handled by a separate path).

Link extraction is, of course, a very important form of document processing. Extensibility in this regard means that it should be easy to add modules to handle Javascript, Flash, and other media-types that have embedded links.

Another form of document processing is duplicate elimination. It should be possible for the crawler to recognize, during the crawl, that a duplicate has been encountered so that special (more efficient) steps can be taken. Any special duplicate handling should be noted in the crawler output and amenable to analysis by standard tools. It should be possible to vary the definition of duplicate used (from very strict to very loose).

Output is another important form of document processing. The output of the crawler should be highly configurable in both the selection of data to be output and the format in which it is output. While one typically saves the pages from a crawl, crawls are often run simply to collect features or statistics. In these cases, the pages themselves are not stored, rather the features and/or statistics are extracted in the crawler and stored directly to disk, reducing disk-space requirements significantly. The crawler must support this type of output.

File formats for output vary. In the case of pages, for example, people use ZIP files, tar files, or the Internet Archives ARC file. We need to support ARC files, whose format we will document, expand where necessary, and promote as a general standard. Our framework must accommodate third-party programmatic extension to other formats when necessary.

#### 3.1.1.6 Interactivity

The Web is becoming increasingly interactive: to receive Web pages, surfers must often fill out forms, and browsers must support authentication and cookie mechanisms.

To keep pace with this ever-more interactive Web, crawlers must support these interaction mechanisms. At the very least, they must support cookies and simple password schemes. And there is an increasing desire for even more support, thus there needs to be extensibility mechanisms to allow crawl-designers to extend the repertoire of interactions supported by the crawl.

### 3.1.2 Extensibility mechanisms

There are two well-known mechanisms for making software extensible: parameterization and code extensions. Parameterization is easy to use, but is restrictive in what it can support. Code extensions requires programming, but supports a wider range of flexibility. As discussed below, the crawler framework should support both forms of extensibility.

In addition to supporting both forms of extensibility, the crawler should support dynamic extension; that is, changing of the crawl while a crawl is in process. We discuss this issue below as well.

#### 3.1.2.1 Parameterization

Extensibility through parameterization is the ability to customize the behavior of the crawl through the use of code modules whose behavior can be changed by supplying appropriate parameters. For example, a frontier module might have a politeness parameter controlling how much time to wait before revisiting a given server, or a scheduling parameter controlling whether image objects from a page are put at the front or the end of the crawling queue. As another example, a filter module might take a set of regular expressions that control what newly-discovered URLs are to be allowed into the frontier.

Another form of parameterization is the ability to control, from a simple, run-time configuration file, which code modules are used during the crawl. For example, one implementation of the frontier might be based on a breadth-first crawling policy; a different implementation might be based on a page-rank-ordered policy. Through simple directives in the crawlers configuration file, the crawl designer should be able to select from among such frontier implementations. (Of course, these modules themselves might be further parameterized; as discussed above, for example, both of these modules might have a flag indicating that images from a page are to be treated as exceptions to the general scheduling policy, ensuring that such images get downloaded contemporaneously with the pages referring to them.)

At a minimum, the crawler framework should have parameters controlling the following aspects of the crawl:

 **Filtering.** Simple URL-based filtering in which patterns specify hostnames and/or paths of interest. Also, it must be possible to filter based on hop counts, i.e., based on the number of hops from a page allowed by the URL-based tests.

 **Scheduling.** As mentioned above, scheduling is determined in large part by the types of frontiers that are allowed by the crawler. We believe the following types of frontiers are necessary:

 **Multiple-priority, breadth-first frontier.** These frontiers allow efficient, disk-based implementations that can scale to very large crawls. Multiple priorities means there are multiple queues with different priorities, allowing one to schedule urgent downloads above the basic queue. Such urgent downloads would support, for example, the download of an image contemporaneously with the page containing it. These priorities can also be used in continuous crawls, wherein frequently-changing pages can be given a higher priority than slowly-changing ones.

 **Site-first frontier.** These frontiers are closer to depth-first crawlers, attempting to complete sites as quickly as possible before starting on new sites. Like breadth-first frontiers, these are very scalable and can support a priority scheme.

 **Dynamically ranked frontier.** The previous frontiers assign an unchanging priority to URLs when they are (re)inserted into the frontier. A dynamically ranked frontier allows the crawl to update URL priorities after insertion. Such frontiers can be used, for example, to download pages in close-to-Pagerank order and also allow ranks to be modified according to linguistic or other considerations. These frontiers typically have an in-memory component that limits their scalability according to the amount of memory on the crawler machines, and thus are less scalable than the statically ranked frontiers.

All of these frontiers must support distribution across multiple machines. They should support continuous crawling and flexible politeness policies.

 **HTTP.** The crawler must come with built-in support for HTTP (v. 1.1). This HTTP implementation should support robot exclusions, but it should also be possible to ignore robot exclusions. This module should support parameterization of outgoing headers (e.g., allowing the crawl designer to set User-Agent and other header fields). This module should support use of HTTPs if-modified-since header (and also support _not_ using it).

 **Link extraction.** The crawler should be able to extract links from HTML and from Flash objects. On HTML objects, it should make a best effort at extracting links from Javascript and other scripting languages.

 **Output.** As mentioned above, output to ARC files, TAR files, and ZIP files needs to be supported. There must be a mechanism for controlling what objects get saved, based both on media-type and also on URL-based tests.

Crawl-designers will not have to write new pieces of code to achieve the above-listed changes.

#### 3.1.2.2 Code extensions

Extensibility through code extensions is the ability to change the behavior of the crawl by inserting into the crawling software customize modules that replace or augment standardized modules. This allows the crawl design to develop (or have developed) modules that meet needs that cannot be met through parameterization of existing modules.

Support for code extensions is difficult. It is one thing for the crawlers original authors write two or three modules from which the crawl operator can choose. It is quite another to design and document appropriate internal APIs that would allow third parties to write their own modules. Such an extensible design cannot be added to a crawler after it has been written; it influences all aspects of the crawlers internal composition and must be contemplated from the very beginning.

An extensible crawling framework must support code extensions in its downloading modules (which impacts extensibility of protocols and interactivity), its frontier (impacts scheduling and politeness), frontier insertion (selection), and document processing (feature extraction, output, discovery, filtering). There must also be a good internal infrastructure for the various modules to communicate their results to one another; this is particularly the case for communicating the results of document processing throughout the crawler.

#### 3.1.2.3 Dynamic extensibility

Crawls often extend over days and even weeks; continuous crawls can run for months. It is not enough that a long-running crawl be extensible when it is first started, it must be possible to change the crawl _during_ the crawl itself. To change policies, it may be necessary to temporarily halt the crawler, but the restart process must be reasonably fast and must not lose track of the state of the crawl.

### 3.2 High bandwidth

A single (adequately-configured) Pentium-based machine should be able to download 10-20M documents per day (which translates roughly to 15-30 Mbps bandwidth averaged over the entire day). Further, it should be possible to distribute a crawl over 10-20 such machines to achieve 10-20x improvements in crawling bandwidth.

### 3.3 Sustained bandwidth

Bandwidth of the crawler should not stop dropping until at least 2.5B documents have been downloaded (if at all). (This assumes no interference from politeness issues, e.g., an adequate pool of sites to pick from.)

### 3.4 Incrementality

The crawler should be usable for small as well as large crawls. While it should be easy to use for all crawls, it should be brain-dead simple for small crawls, i.e., the operator of a small crawl should not be forced to suffer unneeded complexity just because we also want to support large crawls. We believe the best way to support this requirement is a default configuration of the toolkit that does an adequate job for simple crawling problems.

### 3.5 Portability

The crawler must run very well on Linux. Where practical, the crawler should also be capable of small and test runs on Windows 2K/XP. This means it is written in languages and with middleware and libraries that are portable.

### 3.6 Ease of operation

As suggested in 3.4, ease of operation is a somewhat relative term. A simple, single-machine, multi-hour crawl should be attainable by operators with relatively primitive OS administrative skills and no programming skills. A complex, multi-machine, multi-week crawl will require more advanced skills. But in either case, through a combination of simplicity and good documentation, the crawler should be operable by anybody with the requisite OS and programming skills; appeals to a secret oral tradition (unfortunately common when it comes to crawlers) should absolutely not be necessary.

In addition to being simple to operate, the crawler should require little attention or effort while it is underway. While exceptional conditions requiring operator intervention will always exist, the crawler should be as robust to networking and other problems as possible. The operator should not be required to supply an ongoing list of URLs to download or otherwise feed the crawler to keep it going (although it _should_ be possible for an operator to insert URLs if they want to).

Long (multi-day crawls) can be temporarily halted for a number of reasons: the WAN connection becomes unavailable, a bug in the software is tickled, or the crawling policy needs to be changed. To support such interruptions, it must be possible to _restart_ the crawl from data stored on disk. Ideally, such restarts should be clean in that data is not lost; e.g., if a document is being downloaded but does not quite finish when the crawler is shut down, then the crawler should attempt to download that document again when it restarts.

Good logging (including control over log output) is also necessary for monitoring the crawlers behavior and diagnosing any problems. Logging should also support long-term confidence in the output of crawling activity, assisting in the determination of data provenance and authenticity. However, such logs should never become necessary to usefully analyze crawl results: standalone ARC files should remain sufficient and complete as a self-contained crawl record.

Some kind of control panel is needed to allow the operator perform basic operations such as starting a crawl, monitoring the crawls progress, initiating a temporary halt to the crawl or changing the logging parameters. A web-based interface for performing these tasks, from the local machine or remotely, should be offered.

Finally, ease of operation extends to software required by the crawler as well as the crawler itself. For example, although object-request brokers (ORBs) _per se_ are acceptable, many of them are very difficult to administer, and the crawl must not depend on such a one.
