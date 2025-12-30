/**
  @param {number} number
  @param {number} length
*/
function formatNumber(number, length) {
  return number.toString().padStart(length, '0');
}

/**
  Immutable calendar date without a time component.
  Months are 1-indexed.
*/
export class PlainDate {
  #year;
  #month;
  #day;

  /**
    @param {number} year
    @param {number} month
    @param {number} day
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

  toString() {
    const y = formatNumber(this.#year, 4);
    const m = formatNumber(this.#month, 2);
    const d = formatNumber(this.#day, 2);
    return `${y}-${m}-${d}`;
  }

  toJSON() {
    return this.toString();
  }

  /** @param {Date} date */
  static fromDate(date) {
    return new PlainDate(
      date.getFullYear(),
      date.getMonth() + 1, // Date.getMonth() is zero-indexed
      date.getDate()
    );
  }

  toUTC() {
    return Date.UTC(this.#year, this.#month - 1, this.#day)
  }

  /**
    @param {PlainDate} d1
    @param {PlainDate} d2
  */
  static daysBetween(d1, d2) {
    const date1 = d1.toUTC();
    const date2 = d2.toUTC();

    const msPerDay = 24 * 60 * 60 * 1000;
    return Math.round((date2 - date1) / msPerDay);
  }
}

/** Immutable time without a date component */
export class PlainTime {
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

  toJSON() {
    const hh = formatNumber(this.#hour, 2);
    const mm = formatNumber(this.#minute, 2);
    const ss = formatNumber(this.#minute, 2);
    return `${hh}:${mm}:${ss}`;
  }
}

/** @typedef {{
  date: PlainDate;
  time?: PlainTime;
}} FloatingDate */

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

/**
  @param {string} str
  @returns {FloatingDate}
*/
export function parseDate(str) {
  if (str.length > 10) {
    const [year, month, day, hour, minute, second] = str.split(/[- :]/).map(Number);
    return {
      date: new PlainDate(year, month, day),
      time: new PlainTime(hour, minute, second),
    }
  } else {
    const [year, month, day] = str.split('-').map(Number);
    return {
      date: new PlainDate(year, month, day),
    }
  }
}

/** @param {PlainDate} date */
export function formatDate(date) {
  const diffDays = PlainDate.daysBetween(
    PlainDate.fromDate(new Date()),
    date
  );

  const dateUTC = new Date(date.toUTC());

  if (diffDays === 0) {
    return "Today";
  } else if (diffDays === 1) {
    return "Tomorrow";
  } else if (diffDays > 0 && diffDays <= 7) {
    return new Intl.DateTimeFormat("en-US", { weekday: "long" }).format(dateUTC);
  } else {
    return new Intl.DateTimeFormat("en-GB", {
      day: "numeric",
      month: "short",
      year: "numeric",
    }).format(dateUTC);
  }
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

/** @param {number} code */
function isLetter(code) {
  return (code >= 97 && code <= 122) || (code >= 65 && code <= 90);
}

/** @param {number} code */
function isNumber(code) {
  return code >= 49 && code <= 57;
}

/** @typedef {{
  start: number;
  end: number;
}} Range */

/** @typedef { Range & {
  kind: number;
  value: number;
}} Token */

/**
  @param {string} text
  @yields {Token}
*/
function* tokens(text) {
  const { length } = text;

  for (let i = 0; i < length; i++) {
    const start = i;
    const code = text.charCodeAt(i);

    if (isLetter(code)) {
      // Advance until space or end (consumes space)
      while (i < length && text.charCodeAt(i) !== 32) {
        i++;
      }

      const id = keywords.get(
        text.slice(start,i)
          .toLowerCase()
      );

      if (id === undefined) {
        yield {
          kind: 0, // Text
          // no value
          start,
          end: i,
        };
      } else {
        yield {
          kind: 1, // Keyword
          value: id,
          start,
          end: i,
        };
      }

    } else if (isNumber(code)) {
      while (i < length && isNumber(text.charCodeAt(i))) {
        i++;
      }

      yield {
        kind: 2,
        value: parseInt(text.slice(start, i)),
        start,
        end: i,
      };

      // un-consume non-numeric code point
      i--;
    }
  }
}

/** @typedef { Range & FloatingDate } DatetimeExpr */

/**
  Prematurely-optimised iterative LL(1) parser
  @param {string} text
*/
export function parseTaskName(text) {
  const iterator = tokens(text);

  /** @type {boolean} */
  let done;

  /** @type {Token}
    Can persist across loop iterations for 1 lookahead
  */
  let token;

  /**
    State for whether to start a new Expr
    0: nothing
    1: date
    2: date + time
  */
  let block = 0;

  /** @type {DatetimeExpr[]}  */
  const nodes = [];

  /** @type {DatetimeExpr} */
  let node = {};

  function nextToken() {
    ({done, value: token} = iterator.next());
  }

  function reset() {
    if (block > 0) {
      nodes.push(node);
      node = {};
      block = 0;
    }
  }

  /**
    @param {PlainDate} date
  */
  function setDate(plainDate, start = token.start) {
    if (block === 1) {
      reset();
    }

    if (block === 0) {
      delete node.time;
      node.start = start;

      block = 1;
    }

    node.end = token.end;
    node.date = plainDate;
  }

  /**
    @param {PlainTime} plainTime
  */
  function setTime(plainTime, start = token.start) {
    if (block === 2) {
      reset();
    }

    if (block === 0) {
      // default to today
      node.date = PlainDate.fromDate(new Date());
      node.start = start;

      block = 2;
    }
    node.end = token.end;
    node.time = plainTime;
  }

  nextToken(); // Get first token

  topLoop: while (!done) {
    switch (token.kind) {
    // Text
    case 0:
      reset();
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
          const {start} = token;
          nextToken();
          if (done) break topLoop;

          switch (token.value) {
          // week
          case 11:
            // todo: allow configuring day
            setDate(nextDay(1), start);
            break;

          // month
          case 12:
            setDate(inMonths(1), start);
            break;

          // year
          case 13:
            setDate(inYears(1), start);
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
      const {start} = token;
      nextToken();
      if (done) break topLoop;
      if (token.kind !== 1) continue topLoop;

      // Keyword
      switch (token.value) {
      // am
      case 14:
        setTime(new PlainTime(hour, 0, 0), start);
        break;

      // pm
      case 15:
        setTime(new PlainTime(hour + 12, 0, 0), start);
        break;
      }
      break;
    }

    nextToken();
  }

  reset();
  return nodes;
}

// const test = "example task today ";
// for (const t of tokens(test)) {
//   console.log(t);
// }

// console.log(parseTaskName(test));
