export namespace main {
	
	export enum MsgCode {
	    OK = 0,
	    Notfound = 1,
	    NotLogin = 2,
	    IsLogined = 3,
	    VIPExpired = 4,
	    NotSelectGame = 5,
	    Accelerating = 6,
	    InvalidMonths = 7,
	    GameExist = 8,
	    NotAccelerated = 9,
	    RequireGameId = 10,
	    Unknown = 11,
	}
	export class Message {
	    code: MsgCode;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class Stats {
	    stamp: number;
	    gateway_loc: string;
	    forward_loc: string;
	    server_loc: string;
	    rtt_gateway: number;
	    rtt_forward: number;
	    loss_client_uplink: number;
	    loss_client_downlink: number;
	    loss_gateway_uplink: number;
	    loss_gateway_downlink: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stamp = source["stamp"];
	        this.gateway_loc = source["gateway_loc"];
	        this.forward_loc = source["forward_loc"];
	        this.server_loc = source["server_loc"];
	        this.rtt_gateway = source["rtt_gateway"];
	        this.rtt_forward = source["rtt_forward"];
	        this.loss_client_uplink = source["loss_client_uplink"];
	        this.loss_client_downlink = source["loss_client_downlink"];
	        this.loss_gateway_uplink = source["loss_gateway_uplink"];
	        this.loss_gateway_downlink = source["loss_gateway_downlink"];
	    }
	}
	export class StatsList {
	    list: Stats[];
	
	    static createFrom(source: any = {}) {
	        return new StatsList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.list = this.convertValues(source["list"], Stats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

