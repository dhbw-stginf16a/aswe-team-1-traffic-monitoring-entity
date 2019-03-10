# Concern Definition (Traffic)

**Concern-Tag :** `traffic`

Monitoring Entity to get traffic & public transportation information from Google Map's API.

## General Parameters

- **location** (string): Current location of the user
- **destination** (string): Desired destination of the user
- **preference** (string/array): User's prefered method of transportation.
  - _car_
  - _public_
  - _bike_
  - _walk_

## Request Types
### Route to Destination
**Type-Tag:** `traffic_route`

### Issues to Destination
**Type-Tag:** `traffic_issues`

#### Request

- **location**: [see](#general-parameters)
- **destinaten**: [see](#general-parameters)
- **preference**: [see](#general-parameters)

#### Response


#### Example

Request

```json
[Example for Request]
```

Response

```json
[Example for Response]
```

[List of all supported subscription Types]:

## Subscription Types

### [Type Name]

**Type-Tag:** `[Type Tag]`

#### Message

- **[Parameter1 Name]** ([Parameter1 Type]): [Description]
  - _[Possible Value]_: [Value Description]
- **[Parameter2 Name]**: ([Parameter2 Type]): [Description]

#### Example

```json
[Example for message]
```

---

Example Definition

# Concern Definition (Weather)

**Concern-Tag :** `weather`

This concern is all about the weather. Additional info should go here.

## General Parameters

- **condition** (string): Describes the weather condition
  - _sun_: The sun is out
  - _rain_: It is raining

## Request Types

### Current Weather

**Type-Tag:** `weather_current`

#### Request

- **location** (string): Name of the city to check the weather for

#### Response

- **condition**: [see](#general-parameters)

#### Example

Request

```json
{
  "type": "weather_current",
  "payload": {
    "location": "Stuttgart"
  }
}
```

Response

```json
{
  "type": "weather_current",
  "payload": {
    "condition": "sun"
  }
}
```

## Subscription Types

### Weather Changed

**Type-Tag:**: `weather_changed`

#### Message

- **location** (string): Name of the city to check the weather for
- **condition**: [see](#general-parameters)
- **old-condition**: _see condition_

#### Example

```json
{
  "type": "weather_changed",
  "payload": {
    "location": "Stuttgart",
    "condition": "sun",
    "old-condition": "rain"
  }
}
```
