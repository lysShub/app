export namespace main {
	
	export enum MsgCode {
	    OK = 0,
	    Notfound = 1,
	    NotLogin = 2,
	    IsLogined = 3,
	    VIPExpired = 4,
	    NotSetGame = 5,
	    Accelerating = 6,
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
	    gateway_location: string;
	    destinatio_location: string;
	    gameserver_location: string;
	    loss_uplink1: number;
	    loss_downlink1: number;
	    loss_uplink2: number;
	    loss_downlink2: number;
	    ping1: number;
	    ping2: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gateway_location = source["gateway_location"];
	        this.destinatio_location = source["destinatio_location"];
	        this.gameserver_location = source["gameserver_location"];
	        this.loss_uplink1 = source["loss_uplink1"];
	        this.loss_downlink1 = source["loss_downlink1"];
	        this.loss_uplink2 = source["loss_uplink2"];
	        this.loss_downlink2 = source["loss_downlink2"];
	        this.ping1 = source["ping1"];
	        this.ping2 = source["ping2"];
	    }
	}

}

