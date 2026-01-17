local MAX_MONTH = 1;
local MAX_DAY = 14;

local months = std.map(function(n) if n < 10
  then "0" + std.toString(n)
  else std.toString(n), std.range(1, MAX_MONTH));

local days = function (lim) 
  std.map(function(n) if n < 10
    then "0" + std.toString(n)
    else std.toString(n), std.range(1, lim));

local md = function (m) std.map(function (d) "2026-" + m + "-" + d, days(30));

local action = function (ymd) {
    action: "combine",
    name: "combine-summaries",
    table: "itslog_summary",
    source: ymd,
    destination: "summary"};


// A specialized filter to make sure we don't go past a given date.
local filter = function (ls, mlim, dlim) std.filter(function (s) 
  local month = std.parseInt(std.split(s, "-")[1]);
  local day = std.parseInt(std.split(s, "-")[2]);
  if month > mlim 
  then 
    false
  else
    if month == mlim
    then
      if day <= dlim
      then true
      else false
    else true  
  , ls);

{
  server: {
    url: "https://localhost:8443/v1/etl",
  },
  actions: [
    // messages take one param
    // everything can have a "message" option
    {
      action: "message",
      message: "HELO",
    },
    
  ] + std.map(action, filter(std.flatMap(md, months), MAX_MONTH, MAX_DAY))
}
