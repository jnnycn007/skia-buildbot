// DO NOT EDIT. This file is automatically generated.

export interface ValuePercent {
	value: string;
	percent: number;
}

export interface CountHandlerRequest {
	q: string;
	begin: number;
	end: number;
}

export interface CountHandlerResponse {
	count: number;
	paramset: ParamSet;
}

export interface CommitID {
	offset: number;
}

export interface CommitDetailsRequest {
	cid: CommitID;
	traceid: string;
}

export type ParamSet = { [key: string]: string[] };

export type Status = "" | "positive" | "negative" | "untriaged";
