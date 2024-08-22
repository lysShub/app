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
	export class UserInfo {
	    name: string;
	    password: string;
	    phone: string;
	    expire: number;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.password = source["password"];
	        this.phone = source["phone"];
	        this.expire = source["expire"];
	    }
	}
	export class ApiResponse {
	    code: MsgCode;
	    msg: string;
	    data: UserInfo;
	
	    static createFrom(source: any = {}) {
	        return new ApiResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = this.convertValues(source["data"], UserInfo);
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
	export class GameInfo {
	    game_id: number;
	    name: string;
	    icon_path: string;
	    bgimg_path: string;
	    game_servers: string[];
	    cache_game_server: string;
	    cache_fix_route: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GameInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.game_id = source["game_id"];
	        this.name = source["name"];
	        this.icon_path = source["icon_path"];
	        this.bgimg_path = source["bgimg_path"];
	        this.game_servers = source["game_servers"];
	        this.cache_game_server = source["cache_game_server"];
	        this.cache_fix_route = source["cache_fix_route"];
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

