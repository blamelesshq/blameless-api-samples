export namespace App {
    export interface APMResults {
        start: number;
        end: number;
        goodValue?: number;
        validValue?: number;
        latency?: number;
    }

    export interface SLI {
        orgId: number;
        id: number;
        createdAt: number;
        updatedAt: number;
        deletedAt: number;
        name: string;
        description: string;
        dataSourceId: number;
        sliTypeId: number;
        serviceId: number;
        metricPath: string;
        userId: number;
        checkpoint: number;
        backfillJobsTotal: number;
        backfillJobsDone: number;
        backfillJobsFailed: number;
        gcpSettingsId: string;
    }

    export interface SLIRawData {
        id: string;
        sliId: number;
        latency?: number;
        throughput?: number;
        correctness?: number;
        saturation?: number;
        durability?: number;
        validRequest?: number;
        goodRequest?: number;
        start: number;
        end: number;
    }
}
