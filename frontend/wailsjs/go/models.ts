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
	export class ApiResponse {
	    code: MsgCode;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new ApiResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class Loss {
	    gateway: number;
	    forward: number;
	
	    static createFrom(source: any = {}) {
	        return new Loss(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gateway = source["gateway"];
	        this.forward = source["forward"];
	    }
	}
	export class Message {
	    code: MsgCode;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	    }
	}
	export class Stats {
	    stamp: number;
	    gateway: string;
	    forward: string;
	    server: string;
	    ping_gateway: number;
	    ping_forward: number;
	    uplink: Loss;
	    donwlink: Loss;
	    total: Loss;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stamp = source["stamp"];
	        this.gateway = source["gateway"];
	        this.forward = source["forward"];
	        this.server = source["server"];
	        this.ping_gateway = source["ping_gateway"];
	        this.ping_forward = source["ping_forward"];
	        this.uplink = this.convertValues(source["uplink"], Loss);
	        this.donwlink = this.convertValues(source["donwlink"], Loss);
	        this.total = this.convertValues(source["total"], Loss);
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

