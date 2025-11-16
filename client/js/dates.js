/** Immutable calendar date without a time component */
class PlainDate {
  #year;
  #month;
  #day;

  /**
    @arg {number} year
    @arg {number} month
    @arg {number} day
  */
  constructor(year, month, day) {
    this.#day = day;
    this.#month = month;
    this.#year = year;
  }

  get day() {
    return this.#day;
  }

  get month() {
    return this.#month;
  }

  get year() {
    return this.#year;
  }

  /** @arg {Date} date */
  static fromDate(date) {
    return new PlainDate(
      date.getFullYear(),
      date.getMonth() + 1, // Date.getMonth() is zero-indexed
      date.getDate()
    );
  }
}

/** Immutable time without a date component */
class PlainTime {
  #hour;
  #minute;
  #second;

  /**
    @arg {number} hour
    @arg {number} minute
    @arg {number} second
  */
  constructor(hour, minute, second) {
    this.#hour = hour;
    this.#minute = minute;
    this.#second = second;
  }

  get second() {
    return this.#second;
  }

  get minute() {
    return this.#minute;
  }

  get hour() {
    return this.#hour;
  }
}

function nextDay(day) {
  const date = new Date();
  let daysUntil = (day - date.getDay() + 7) % 7;

  if (daysUntil === 0) daysUntil = 7;

  date.setDate(date.getDate() + daysUntil);
  return PlainDate.fromDate(date);
}

function inDays(n) {
  const date = new Date();
  date.setDate(date.getDate() + n);
  return PlainDate.fromDate(date);
}

function inMonths(n) {
  const date = new Date();
  date.setMonth(date.getMonth() + n);
  return PlainDate.fromDate(date);
}

function inYears(n) {
  const date = new Date();
  date.setFullYear(date.getFullYear() + n);
  return PlainDate.fromDate(date);
}

const keywords = new Map([
  ["sunday", 0],
  ["monday", 1],
  ["tuesday", 2],
  ["wednesday", 3],
  ["thursday", 4],
  ["friday", 5],
  ["saturday", 6],

  ["today", 7],
  ["tomorrow", 8],
  ["someday", 9],

  ["next", 10],
  ["week", 11],
  ["month", 12],
  ["year", 13],

  ["am", 14],
  ["pm", 15]
]);

/** @arg {number} code */
function isLetter(code) {
  return (code >= 97 && code <= 122) || (code >= 65 && code <= 90);
}

/** @arg {number} code */
function isNumber(code) {
  return code >= 49 && code <= 57;
}

/** @typedef {{
  kind: number;
  value: number;
}} Token */

/**
  @arg {string} text
  @yields {Token}
*/
function* tokens(text) {
  const { length } = text;

  for (let i = 0; i < length; i++) {
    const start = i;
    const code = text.charCodeAt(i);

    if (isLetter(code)) {
      // Advance until space or end
      while (i < length && text.charCodeAt(i) !== 32) {
        i++;
      }

      const id = keywords.get(
        text.slice(start,i)
          .toLowerCase()
      );

      if (id === undefined) {
        // Text
        yield { kind: 0 };
      } else {
        // Keyword
        yield { kind: 1, value: id};
      }

    } else if (isNumber(code)) {
      while (i < length && isNumber(text.charCodeAt(i))) {
        i++;
      }

      yield {
        kind: 2,
        value: parseInt(text.slice(start, i))
      };

      // un-consume non-numeric code point
      i--;
    }
  }
}

/** @arg {string} text */
function parse(text) {
  const iterator = tokens(text);

  /** @type {boolean} */ let done;
  /** @type {Token} */ let token;

  const nextToken = () => {
    ({done, value: token} = iterator.next());
  };
  nextToken(); // Get first token

  /** State for whether to set a new date/time */
  let inDatetimeBlock = false;

  /** @type {PlainDate} */ let date;
  /** @type {PlainTime} */ let time;

  /** @arg {PlainDate} date */
  function setDate(plainDate) {
    if (!inDatetimeBlock) {
      time = undefined;
      inDatetimeBlock = true;
    }
    date = plainDate;
  }

  /** @arg {PlainTime} plainTime */
  function setTime(plainTime) {
    if (!inDatetimeBlock) {
      // default to today
      date = PlainDate.fromDate(new Date());
      inDatetimeBlock = true;
    }
    time = plainTime;
  }

  topLoop:
  while (!done) {
    switch (token.kind) {
    // Text
    case 0:
      inDatetimeBlock = false;
      break;

    // Keyword
    case 1:
      if (token.value < 7) {
        setDate(nextDay(token.value));
      } else {
        switch (token.value) {
        // today
        case 7:
          setDate(PlainDate.fromDate(new Date()));
          break;

        // tomorrow
        case 8:
          setDate(inDays(1));
          break;

        // someday
        case 9:
          setDate(inMonths(2));
          break;

        // next
        case 10:
          nextToken();
          if (done) break topLoop;

          switch (token.value) {
          // week
          case 11:
            // todo: allow configuring day
            setDate(nextDay(1));
            break;

          // month
          case 12:
            setDate(inMonths(1));
            break;

          // year
          case 13:
            setDate(inYears(1));
            break;

          default:
            continue topLoop;
          }
          break;
        }
      }
      break;

    // Number
    case 2:
      const hour = token.value;
      nextToken();
      if (done) break topLoop;
      if (token.kind !== 1) continue topLoop;

      // Keyword
      switch (token.value) {
      // am
      case 14:
        setTime(new PlainTime(hour));
        break;

      // pm
      case 15:
        setTime(new PlainTime(hour + 12, 0, 0));
        break;
      }
      break;
    }

    nextToken();
  }

  // return result
  console.log(date?.day, time?.hour);
}

const test = "tomorrow example 6am";
// for (const t of tokens(test)) {
//   console.log(t);
// }
parse(test);
