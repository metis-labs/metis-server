# Contribution Guide

Metis accepts contributions via GitHub pull requests. This document outlines some of the conventions on commit message formatting, contact points for developers, and other resources to help get contributions into Metis.

## Getting Started

- Fork the repository on GitHub
- Read the README.md for build instructions

## Contributor's Workflow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where to base the contribution. This is usually master
- Make commits of logical units
- Make sure commit messages are in the proper format
- Push changes in a topic branch to a personal fork of the repository
- Submit a pull request to metis/metis-server
- The PR must receive a LGTM from maintainers

## Coding Style

The coding style suggested by the Golang community is used in Metis. See the [style doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

## Commit Message

We follow a rough convention for commit messages that is designed to answer two questions: what changed and why. The subject line should feature the what and the body of the commit should describe the why.

```
Remove the synced seq when detaching the document

To collect garbage like CRDT tombstones left on the document, all
the changes should be applied to other replicas before GC. For this
, if the document is no longer used by this client, it should be
detached.
```

The first line is the subject and should be no longer than 70 characters, the second line is always blank, and other lines should be wrapped at 80 characters. This allows the message to be easier to read on GitHub as well as in various git tools.
