/**
 * Gets temperature data for the last N days averaged for each day.
 */
function getTempData(projectId, deviceId, days) {
    response = _runQuery(projectId,
        "SELECT" +
        "  PARSE_UTC_USEC(STRFTIME_UTC_USEC(timestamp, '%Y-%m-%d %H:00:00')) / 1000 time, " +
        "  AVG(humidity) humidity, " +
        "  AVG(temp) temp " +
        "FROM " +
        "  [" + projectId + ":homesensors.sensordata] " +
        "WHERE " +
        "  deviceid = \"" + deviceId + "\" AND " +
        "  timestamp > DATE_ADD(CURRENT_TIMESTAMP(), -" + days + ", 'DAY') " +
        "GROUP BY time " +
        "ORDER BY time;"
    );
    return response;
}

/**
 * Run a BigQuery Query. Walks through all pages of the results and
 * builds the result data. Results are return as a list of objects
 * with fields equal to the fields selected in the query.
 */
function _runQuery(projectId, query) {
    var request = {query: query};

    var queryResults = BigQuery.Jobs.query(request, projectId);
    var jobId = queryResults.jobReference.jobId;

    // Check on status of the Query Job.
    var sleepTimeMs = 500;
    while (!queryResults.jobComplete) {
        Utilities.sleep(sleepTimeMs);
        sleepTimeMs *= 2;
        queryResults = BigQuery.Jobs.getQueryResults(projectId, jobId);
    }

    // Get all the rows of results.
    var rows = queryResults.rows;
    while (queryResults.pageToken) {
        queryResults = BigQuery.Jobs.getQueryResults(projectId, jobId, {
            pageToken: queryResults.pageToken
        });
        rows = rows.concat(queryResults.rows);
    }

    // Parse the result rows into the final list.
    results = [];
    for (var i=0; i < rows.length; i++) {
        var result = {}; 
        for (var j=0; j < queryResults.schema.fields.length; j++) {
            var field = queryResults.schema.fields[j]; 
            var value = rows[i].f[j].v;
            if (field.type === "INTEGER") {
                value = parseInt(value, 10);
            } else if (field.type === "FLOAT") {
                value = parseFloat(value);
            } else if (field.type === "BOOLEAN") {
                value = value && value == "true";
            } // TODO: TIMESTAMP and RECORD
            result[field.name] = value;
        }
        results.push(result);
    }

    return results;
}

function doGet(e) {
   var tmpl = HtmlService.createTemplateFromFile('Graph');
   tmpl.projectId = e.parameter.project;
   tmpl.deviceId = e.parameter.device;
   return tmpl.evaluate().setSandboxMode(HtmlService.SandboxMode.IFRAME);
}
