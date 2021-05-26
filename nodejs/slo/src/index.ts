import readlineSync from 'readline-sync';
import fetch from 'node-fetch';
import dotenv from 'dotenv';
import sliSelect from 'cli-select';
import { spinner } from './utils';
import { App } from './interfaces/app';

dotenv.config();

let tokenExpiresAtDate = 0;
let token = '';
let tenantId = 0;


/**
 * @description Returns the auth token needed to make calls into the Blameless Public API.
 * @return Promise<string> JWT Token to authorize through Blameless API.
 */
const getAuth0Token = async () => {
    if (tokenExpiresAtDate) {
        const remaining = tokenExpiresAtDate - Date.now();
        if (remaining > 0) {
            return token;
        }
    }

    const response = await fetch('https://blamelesshq.auth0.com/oauth/token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            client_id: process.env.AUTHZERO_CLIENT_ID,
            client_secret: process.env.AUTHZERO_CLIENT_SECRET,
            audience: process.env.AUTHZERO_API_AUDIENCE,
            grant_type: 'client_credentials',
        }),
    });
    const data = await response.json();
    const nowDate = new Date();
    tokenExpiresAtDate = nowDate.setTime(data.expires_in * 1000 + nowDate.getTime());
    token = data.access_token;
    return token;
};

/**
 * @description Returns the tenant id for the host you have under "BLAMELESS_HOST" in your env variables.
 * @return Promise<number> TenantId needed to perform some queries to the Blameless public API.
 */
const getTenantId = async (): Promise<number> => {
    if (tenantId) return tenantId;
    const response = await fetch(`${process.env.BLAMELESS_HOST}/api/v1/identity/tenant`);
    const data = await response.json();
    tenantId = Number(data.tenant_id);
    return tenantId;
};

/**
 * @description Returns the SLI from the Blameless API.
 * @param sliId number - sli ID to fetch from Blameless API.
 * @return Promise<SLI | null> Returns SLI to work with.
 */
const getSliById = async (sliId: number): Promise<App.SLI | null> => {
    const response = await fetch(`${process.env.BLAMELESS_HOST}/api/v1/services/SLOServiceCrud/GetSLI`, {
        method: 'POST',
        headers: {
            'content-type': 'application/json',
            authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
            orgId: tenantId,
            id: sliId,
        }),
    });
    if (response.status !== 200) {
        return null;
    }
    const data = await response.json();
    return data;
};

/**
 * @description Fetch and return Latency data from the robustperception Prometheus demo site.
 * @example Example of data fetched from robustperception demo site.
 *  {
 *      ...
 *      "result": [
 *          {
 *              "metric": {
 *                  ...
 *              },
 *              "values": [
 *                  [1621362566, "536"],
 *                  [1621362626, "538"],
 *                  [1621362686, "541"],
 *              ]
 *          }
 *      ]
 * }
 *
 * @param minutesRangeInterval number - Interval range in minutes.
 * @param step number - Step in seconds.
 * @returns Promise<APMResults[]> - Example {start: number; end: number; latency: number;}
 */
const getLatencyApmData = async (minutesRangeInterval: number, step: number): Promise<App.APMResults[]> => {
    var now = Date.now();
    var from = new Date(now - 1000 * 60 * minutesRangeInterval);
    var to = new Date(now);
    var start = Math.floor(from.getTime() / 1000);
    var end = Math.floor(to.getTime() / 1000);
    var query =
        'alertmanager_http_request_duration_seconds_bucket{handler="/", job="alertmanager", le="0.5", method="get"}';
    const response = await fetch(
        `http://demo.robustperception.io:9090/api/v1/query_range?end=${end}&start=${start}&step=${step}&query=${encodeURIComponent(
            query,
        )}`,
        {
            headers: {
                'content-type': 'application/json',
            },
        },
    );
    const result = await response.json();
    return result.data.result[0].values.map((value: [number, string]) => ({
        start: value[0],
        end: value[0] + step,
        latency: Number(value[1]),
    }));
};

/**
 * @description Fetch and return Availability data from the robustperception Prometheus demo site.
 * @example Example of data fetched from robustperception demo site.
 *  {
 *      ...
 *      "result": [
 *          {
 *              "metric": {
 *                  ...
 *              },
 *              "value": [
 *                  [1621621113, "1.5399384024639016"],
 *              ]
 *          }
 *      ]
 * }
 *
 * @param minutesRangeInterval number - Interval range in minutes.
 * @returns Promise<APMResults[]> - Example {start: number; end: number; goodValue: number; validValue: number;}
 */
const getAvailabilityApmData = async (minutesRangeInterval: number): Promise<App.APMResults[]> => {
    var now = Date.now();
    var from = new Date(now - 1000 * 60 * minutesRangeInterval);
    var to = new Date(now);
    var start = Math.floor(from.getTime() / 1000);
    var end = Math.floor(to.getTime() / 1000);
    var goodQuery = `sum(rate(prometheus_http_requests_total{code="200"}[${minutesRangeInterval}m]))`;
    const goodResponse = await fetch(
        `http://demo.robustperception.io:9090/api/v1/query?time=${end}&query=${encodeURIComponent(goodQuery)}`,
        {
            headers: {
                'content-type': 'application/json',
            },
        },
    );
    var validQuery = `(sum(rate(prometheus_http_requests_total{code="200"}[${minutesRangeInterval}m]))+sum(rate(prometheus_http_requests_total{code="400"}[${minutesRangeInterval}m])))`;
    const validResponse = await fetch(
        `http://demo.robustperception.io:9090/api/v1/query?time=${end}&query=${encodeURIComponent(validQuery)}`,
        {
            headers: {
                'content-type': 'application/json',
            },
        },
    );
    const goodResults = await goodResponse.json();
    const validResults = await validResponse.json();
    const goodValue = Number(goodResults.data.result[0].value[1]);
    const validValue = Number(validResults.data.result[0].value[1]);
    return [
        {
            start,
            end,
            goodValue: Math.floor(goodValue * 100), // multiplied by 100 for demo purposes
            validValue: Math.floor(validValue * 100), // multiplied by 100 for demo purposes
        },
    ];
};

/**
 * @description Takes APM data results, and ingest them into Blameless
 * @param sliType string - latency or availability.
 * @param sliID number - ID of the SLI with want to ingest data on.
 * @param apmValues APMResults[] - APM data, Example {start: number; end: number; goodValue: number; validValue: number;}.
 * @returns Promise<SLIRawData[]> - Example {start: number; end: number; goodValue: number; validValue: number;}
 */
const ingestApmData = async (
    sliType: string,
    sliId: number,
    apmValues: App.APMResults[],
): Promise<App.SLIRawData[]> => {
    let rawDatas;
    if (sliType === 'latency') {
        rawDatas = apmValues.map((data) => ({
            sliId,
            latency: data.latency,
            start: data.start,
            end: data.end,
        }));
    }
    if (sliType === 'availability') {
        rawDatas = apmValues.map((data) => ({
            sliId,
            goodRequest: data.goodValue,
            validRequest: data.validValue,
            start: data.start,
            end: data.end,
        }));
    }

    const response = await fetch(
        `${process.env.BLAMELESS_HOST}/api/v1/services/SLOTimeSeriesService/SliRawDataPostMany`,
        {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
                authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
                orgId: tenantId,
                sliType: 'latency',
                rawDatas,
            }),
        },
    );
    const result = await response.json();
    return result.sliRawData;
};

/**
 * @description Wrapper method, fetch APM data (getLatencyApmData) and then ingest it into Blameless (ingestApmData).
 * @param sliType string - latency or availability.
 * @param sliID number - ID of the SLI with want to ingest data on.
 * @param minutesRangeInterval number - Interval range in minutes.
 * @param step number - Step in seconds.
 */
const getApmDataAndIngest = async (sliType: string, sliId: number, minutesRangeInterval: number, step?: number) => {
    let apmValues;

    // getting APM values
    const apmValuesSpinner = spinner('💹 Getting APM values').start();
    switch (sliType) {
        case 'latency':
            try {
                if (step) {
                    apmValues = await getLatencyApmData(minutesRangeInterval, step);
                    apmValuesSpinner.succeed();
                }
            } catch (error) {
                apmValuesSpinner.fail();
                return console.error('Error: ', error);
            }
            break;
        case 'availability':
            try {
                apmValues = await getAvailabilityApmData(minutesRangeInterval);
                apmValuesSpinner.succeed();
            } catch (error) {
                apmValuesSpinner.fail();
                return console.error('Error: ', error);
            }
            break;
    }

    if (apmValues) {
        // ingesting APM values
        const ingestApmValuesSpinner = spinner('⏬ Ingesting APM values').start();
        try {
            await ingestApmData(sliType, sliId, apmValues);
            ingestApmValuesSpinner.succeed();
        } catch (error) {
            ingestApmValuesSpinner.fail();
            return console.error('Error: ', error);
        }
    }
};

/**
 * @description Helper method to render a spinner between ingestion intervals.
 * @param minutesRangeInterval number - Interval range in minutes.
 */
const showWaitingIntervalSpinner = async (minutesRangeInterval: number) => {
    const intervalSpinner = spinner(`Next ingest will be executed in ${minutesRangeInterval} minute(s)`).start();
    setTimeout(async () => {
        intervalSpinner.stop();
        intervalSpinner.clear();
    }, minutesRangeInterval * 60 * 1000);
};

/**
 * @description Main method to start ingesting data into Blameless API.
 */
const run = async () => {
    try {
        let step: number | undefined;

        // Ask for initial data
        console.info('Select SLI Type (Availability or Latency):');
        const sliTypeResult = await sliSelect({ values: ['availability', 'latency'] });
        const sliType = sliTypeResult.value;
        console.info(`${sliType} selected`);
        const sliId = Number(await readlineSync.question('Insert SLI ID you want to work with:'));
        const minutesRangeInterval = Number(await readlineSync.question('Insert ingest interval (in minutes):'));
        if (sliType === 'latency') {
            step = await Number(readlineSync.question('Insert step (in seconds):'));
            if (!sliId || !minutesRangeInterval || !step || !sliType) {
                return console.error(
                    'Error: ',
                    'You need to insert numeric values for "SLI ID", "Time Range" and "Step"',
                );
            }
        }
        if (sliType === 'availability') {
            if (!sliId || !minutesRangeInterval || !sliType) {
                return console.error('Error: ', 'You need to insert numeric values for "SLI ID" and "Time Range"');
            }
        }

        console.log('\n');

        // auth token
        const authTokenSpinner = spinner('🔒 Getting Auth token').start();
        try {
            await getAuth0Token();
            authTokenSpinner.succeed();
        } catch (error) {
            authTokenSpinner.fail();
            return console.error('Error: ', error);
        }

        // tenant id
        const orgIdSpinner = spinner('🆔 Getting Org Id').start();
        try {
            await getTenantId();
            orgIdSpinner.succeed();
        } catch (error) {
            orgIdSpinner.fail();
            return console.error('Error: ', error);
        }

        // sli data
        const sliIdSpinner = spinner(`👀 Looking for ${sliType} SLI with ID ${sliId}`).start();
        try {
            const sli = await getSliById(sliId);
            if (sli) {
                sliIdSpinner.succeed();
            } else {
                return sliIdSpinner.fail(`Cannot fetch SLI with ID: ${sliId}`);
            }
        } catch (error) {
            sliIdSpinner.fail();
            return console.error('Error: ', error);
        }

        await getApmDataAndIngest(sliType, sliId, minutesRangeInterval, step);
        showWaitingIntervalSpinner(minutesRangeInterval);
        setInterval(async () => {
            await getApmDataAndIngest(sliType, sliId, minutesRangeInterval, step);
            showWaitingIntervalSpinner(minutesRangeInterval);
        }, minutesRangeInterval * 60 * 1000);
    } catch (error) {
        return console.error('Error: ', error);
    }
};

run();
