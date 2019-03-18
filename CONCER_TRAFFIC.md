# Concern Definition (Traffic)

**Concern-Tag :** `traffic`

Monitoring Entity to get traffic & public transportation information from Google Map's API.

## General Parameters

- **location** (string): Current location of the user
- **destination** (string): Desired destination of the user
- **arriveby** (string) [optional]: Desired time of arrival at the specified destination. Seconds since January 1, 1970 UTC
- **travelmode** (array of strings) [optional]: User's preferred method of transportation. Default is just driving.
  - _driving_
  - _transit_ : Public Transport
  - _bicycling_
  - _walking_

## Request Types

### Route to Destination

**Type-Tag:** `traffic_route`

#### Request

- **location**: [see](#general-parameters)
- **destination**: [see](#general-parameters)
- **arriveby**: [see](#general-parameters)
- **travelmode**: [see](#general-parameters)

#### Response

- **routes**: (array): Array of possible routes
  - **travelmode**: single travelmode string [see](#general-parameters)
  - **duration** (number): Estimated time for route in seconds
  - **durationText** (string): Time in text representation
  - **delay** (number): Delay in seconds (Only with travelmode driving).
  - **delayText** (string): Time in text representation
  - **distance** (string): Distance in human readable format
  - **destination** (string): Destination as understood by Google API
  - **location** (string): Location as understood by Google API
  - **link**: (string) [WIP]: Link to Google Maps

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
