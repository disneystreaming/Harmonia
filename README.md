# Harmonia <!-- omit in toc -->

<!-- include link to GH build status here if we like in future -->

___

## Table of Contents <!-- omit in toc -->
- [Overview](#overview)
- [Building and Running :hammer: :running:](#building-and-running-hammer-running)
- [How to Use Harmonia](#how-to-use-harmonia)
  - [What is an RFC?](#what-is-an-rfc)
  - [How do I structure an RFC?](#how-do-i-structure-an-rfc)
  - [Typical Harmonia Workflow](#typical-harmonia-workflow)
    - [Step 1: Create your RFC](#step-1-create-your-rfc)
    - [Step 2: Submit your RFC via `/submitRequest`](#step-2-submit-your-rfc-via-submitrequest)
    - [Step 3: Wait for Stakeholder Responses to come in via `/reviewRequest`](#step-3-wait-for-stakeholder-responses-to-come-in-via-reviewrequest)
    - [Step 4: Analyze Feedback and Submit Updates via `/updateRequest`](#step-4-analyze-feedback-and-submit-updates-via-updaterequest)
    - [Step 5: Wait for Another Round of Stakeholder Responses to come in via `/reviewRequest`](#step-5-wait-for-another-round-of-stakeholder-responses-to-come-in-via-reviewrequest)
    - [Step 6: Get your RFC Accepted into the Schema!](#step-6-get-your-rfc-accepted-into-the-schema)
___

## Overview

> In Greek mythology, Harmonia (/hɑːrˈmoʊniə/; Ancient Greek: Ἁρμονία) is the immortal goddess of harmony and concord.

Harmonia is an API that allows for easy manipulation and review of data schemas via an RFC process that leverages Git
as a backing store.

## Building and Running :hammer: :running:

1. First, go ahead and set up the environment variables depicted below.

Environment Variables
| Variable Name       | Description                                        | Default Value |
| ------------------- | ---------------------------------------------------| ------------- |
| IS_LOCAL            | Set to `true` if you are running the stack locally | `true`        |
| GIT_TOKEN           | Set to GitHub user access token                    | None          |
| GIT_MACHINE_TOKEN   | Set to GitHub machine access token                 | None          |
| TRACKING_REPOSITORY | Set to GitHub tracking repository                  | None          |

For convenience, a script has been provided to set these environment variables locally. Simply run the following to
initialize your local environment.
```
make local && source localenv
```

1. Run `make compile` to compile the source code into the `bin` directory.
2. Run `make run` to run the compiled source.
3. View your changes locally [here](http://localhost:8080)!
4. Considering this is a local build, you must use the `http` scheme in Swagger.

## How to Use Harmonia

This section goes over the fundamentals of how Harmonia should be used to enact schema changes!

### What is an RFC?

An RFC or "Request For Comments" in Harmonia is the payload structure used to submit a proposal for any desired changes
that you may have for the existing schema set. More succinctly, it is a list of actions that you want performed to the
schema set. As the RFC goes through our approval process, it is constantly updated with new actions to track the life of
your request.

### How do I structure an RFC?

At its core an RFC is a list of actions `[{action 1}, {action 2}, {action 3}...]`

Each action has an `actionType`, which must be one of: `add`, `update`, `comment`, `approve` or `load`. As stated above,
each `actionType` is either an action you want performed on the schema OR action metadata on the submitted RFC. For
example, the `add` and `update` action types would be used to `add` and/or `update` a schema entity. But, the `comment`,
`approve` and `load` actions correspond to actions that occurred either by you or others during the lifecycle of the
RFC.

The next piece of the RFC is what the action is acting upon, also known as the `target`. The `target` is an object
that looks like the following:
```
{
  "targetType": ...
  "targetDescriptor": ...
  "lookupKey": ...
  "lookupValue": ...
  "relType": ...
}
```

`targetType` must be one of `item`, `action` or `rfc`. Use `item` if you are acting on a schema entity. Use `action` if
you are acting on an RFC action within the current RFC (this will mainly be used by Harmonia behind the scenes to track
comments). Lastly, use `rfc` if you are acting on this RFC as a whole (again, this will mainly be used by Harmonia
behind the scenes for approval flow).

`targetDescriptor` is optional, **unless** your `targetType` is `item` because Harmonia needs to know which `item`
category you want to work on.

`lookupKey` is the attribute field that Harmonia should use to actually find your target. For example, this could be
an `id`.

`lookupValue` is the value that the `lookupKey` should match to on the target. For example, if you are looking for an
entity with `id` `3` then `lookupKey` should be `id` and `lookupValue` should be `3`.

Lastly, there is a freeform `data` object that can exist at the same level as the `actionType` and `target` properties
of the RFC. This `data` object can contain various properties that are specific to your data schema.

### Typical Harmonia Workflow

Now we will outline a common workflow of taking an RFC from ideation to approval and acceptance into the specification.
One schema change request that occurs often is editing the list of accepted values for a field in our database. Let's
say we want to add `acceptedValue4` to our field named `OurField`.

#### Step 1: Create your RFC

Craft an RFC that tells Harmonia we want to add `acceptedValueFour` to our accepted values for the `OurField` field. Our
RFC would look like this:
```
{
  "actions": [
    {
      "actionType": "update",
      "data": {
        "acceptedValues": "acceptedValue1;acceptedValue2;acceptedValue3;acceptedValueFour"
      },
      "target": {
        "targetDescriptor": "acceptedValueChecker",
        "targetType": "item",
        "lookupKey": "name",
        "lookupValue": "OurField"
      }
    }
  ]
}
```

In the above RFC, we are adding `acceptedValueFour` to the accepted values for `OurField`. As you can see, the `update`
action acts as an **overwrite** on the existing data instead of an addition to it. It is important to keep this in mind
when performing updates. The reason we didn't use an `add` action here is because the `OurField` `acceptedValueChecker`
already exists, we simply want to update it.

#### Step 2: Submit your RFC via `/submitRequest`

Now that we have our RFC, we can submit a POST request to the `/submitRequest` endpoint to officially offer our RFC up
for review.

Now is the time when stakeholders of the `OurField` field will want to weigh in on our request.

#### Step 3: Wait for Stakeholder Responses to come in via `/reviewRequest`

Let's say a stakeholder comes along and thinks that you shouldn't be using the word "Four" inside your accepted
value. Meaning instead he or she wishes our new value was `AcceptedValue4`. He or she would submit the following
payload via the `/reviewRequest` endpoint.

```
{
  "rfcIdentifier": "123456", // this will be known after submitting your initial request
  "topLevelComment": "nice! But you shouldn't use the word "Four" in the value, instead just use "4"",
  "type": "COMMENT"
}
```

If you look at the `/reviewRequest` endpoint you will notice that we also accept an object of `comments`. So why do we
have a `topLevelComment` and a `comments` object? Well in GitHub during a review you can comment on the entire pull
request and individual lines. Our goal is to emulate this functionality. Therefore, use the `topLevelComment` to comment
on the entire RFC and use the `comments` object, with individual action `signatures` as keys and your desired comments
as the values, to comment on a specific action in the RFC.

To give more insight into the `rfc` and `action` target types [described here](#how-do-i-structure-an-rfc), after the
above comment is added the RFC will be updated in the background to include the following action:

```
{
    "actionType": "comment",
    "target": {
        "targetType": "rfc",
        "lookupKey": "signature",
        "lookupValue": "thisisthesignatureoftherfc"
    },
    "data": {
        "comment": "nice! But you shouldn't use the word "Four" in the value, instead just use "4""
    }
}
```

If there were comments on individual actions that were created via the `comments` object then we would see a single
`comment` action like the one above for each comment targeting individual actions by using the `action` `targetType`.

Lastly, there are three types of reviews allowed via this endpoint: `COMMENT`, `REQUEST_CHANGES` and `APPROVE`, which
all correspond directly back to their analogs in GitHub when reviewing a pull request.

#### Step 4: Analyze Feedback and Submit Updates via `/updateRequest`

At this point, you would notice the comment on your RFC and could submit an update to your RFC to match the suggestions.
Note: the way that the `/updateRequest` endpoint works is that it takes whatever you give it, as the new RFC in total.
It doesn't merge or append your update to the existing RFC. Because of this point, our payload would look like this if
we were to address the comment from the stakeholder:

```
{
  "rfc": {
    "actions": [
      {
        "actionType": "update",
        "data": {
          "enum": "acceptedValue1;acceptedValue2;acceptedValue3;acceptedValue4"
        },
        "target": {
          "targetDescriptor": "acceptedValueChecker",
          "targetType": "item",
          "lookupKey": "name",
          "lookupValue": "OurField"
        }
      }
    ]
  },
  "rfcIdentifier": "123456
}
```

You can see we still listed all the accepted values, and if there were other actions in the original RFC that should
still be included we would want to include those in our update RFC above or else they would be **overwritten**!

#### Step 5: Wait for Another Round of Stakeholder Responses to come in via `/reviewRequest`

After submitting the update, the stakeholders could again review. Let's say that everything looks good to them! They
will submit an approval via the `/reviewRequest` endpoint. Their review payload would look like this:

```
{
  "rfcIdentifier": "123456",
  "topLevelComment": "awesome work!,
  "type": "APPROVE"
}
```

#### Step 6: Get your RFC Accepted into the Schema!

Once your RFC has the desired number of approvals it will automatically be integrated into the specification. You can
easily check the status of the loading process of your RFC by using the `/status` endpoint with your assigned
`rfcIdentifier`.