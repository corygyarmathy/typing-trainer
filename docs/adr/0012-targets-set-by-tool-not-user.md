# ADR 0012: Targets set by tool, not by user

- **Status:** Proposed
- **Date:** 2026-06-18

## Context

Improving at typing is a process of building familiarity, from familiarity accuracy, and from accuracy speed. Many tools have mechanisms for managing the accuracy or words-per-minute (WPM) targets, and often suggest for the user to raise these targets over time. But, in my opinion, this manual intervention should be unnecessary. It is a distraction from the core activity, improving at typing, and the suggested procedures for improving should instead be built into the tool itself.

## Decision

The tool will automatically set any given thresholds, be it accuracy, WPM, or anything else. The targets will be based on the user's current performance, and will be highly dynamic - the tool will not hold the user back if they are showing rapid improvement, but it also will not move the user forward it they are not improving.

This follows the existing approach of 'unlocking' keys or ngrams once the user has met the targets of the previous items.

For this to be effective, it will be necessary to provide feedback to the user, such that they know what the targets / thresholds are and why.

## Consequences

**Positive**

- User is able to focus only on improving their typing, allowing the tool to respond to their demonstrated competency and focus where improvement is needed, providing a constant small challenge to get incrementally better.
- Decreased interface, application logic, and database complexity.

**Negative**

- Decreased flexibility for the user, they do not have the ability to manually set targets / thresholds.
