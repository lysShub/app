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
	export class GameInfo {
	    game_id: number;
	    name: string;
	    icon_path: string;
	    bgimg_path: string;
	    game_servers: string[];
	    cache_game_server: string;
	    cache_fix_route: boolean;
	    last_active: number;
	    duration: number;
	    flow: number;
	
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
	        this.last_active = source["last_active"];
	        this.duration = source["duration"];
	        this.flow = source["flow"];
	    }
	}
	export class Message[[]main.GameInfo] {
	    code: MsgCode;
	    msg: string;
	    data: GameInfo[];
	
	    static createFrom(source: any = {}) {
	        return new Message[[]main.GameInfo](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = this.convertValues(source["data"], GameInfo);
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
	export class Message[main.GameInfo] {
	    code: MsgCode;
	    msg: string;
	    data: GameInfo;
	
	    static createFrom(source: any = {}) {
	        return new Message[main.GameInfo](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = this.convertValues(source["data"], GameInfo);
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
	export class Message[main.StatsList] {
	    code: MsgCode;
	    msg: string;
	    data: StatsList;
	
	    static createFrom(source: any = {}) {
	        return new Message[main.StatsList](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = this.convertValues(source["data"], StatsList);
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
	export class Message[main.UserInfo] {
	    code: MsgCode;
	    msg: string;
	    data: UserInfo;
	
	    static createFrom(source: any = {}) {
	        return new Message[main.UserInfo](source);
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
	export class Message[string] {
	    code: MsgCode;
	    msg: string;
	    data: string;
	
	    static createFrom(source: any = {}) {
	        return new Message[string](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class Message[struct {}] {
	    code: MsgCode;
	    msg: string;
	    // Go type: struct {}
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Message[struct {}](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = this.convertValues(source["data"], Object);
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

