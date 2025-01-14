export namespace mysql {
	
	export class MysqlProcess {
	    id: number;
	    user: string;
	    host: string;
	    dbname: string;
	    command: string;
	    time: number;
	    status: string;
	    sql: string;
	    key: string;
	
	    static createFrom(source: any = {}) {
	        return new MysqlProcess(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user = source["user"];
	        this.host = source["host"];
	        this.dbname = source["dbname"];
	        this.command = source["command"];
	        this.time = source["time"];
	        this.status = source["status"];
	        this.sql = source["sql"];
	        this.key = source["key"];
	    }
	}
	export class MysqlProcessF {
	    id: number;
	    user: string;
	    host: string;
	    dbname: string;
	    command: string;
	    time: number;
	    status: string;
	    sql: string;
	    key: string;
	    children: MysqlProcess[];
	
	    static createFrom(source: any = {}) {
	        return new MysqlProcessF(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user = source["user"];
	        this.host = source["host"];
	        this.dbname = source["dbname"];
	        this.command = source["command"];
	        this.time = source["time"];
	        this.status = source["status"];
	        this.sql = source["sql"];
	        this.key = source["key"];
	        this.children = this.convertValues(source["children"], MysqlProcess);
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

export namespace redis {
	
	export class RedisKey {
	    key: string;
	    name: string;
	    type: string;
	    expire: string;
	    size: string;
	
	    static createFrom(source: any = {}) {
	        return new RedisKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.expire = source["expire"];
	        this.size = source["size"];
	    }
	}
	export class RedisMomery {
	    momeryForever: number;
	    countForever: number;
	    momerys3: number;
	    counts3: number;
	    momeryb3s7: number;
	    countb3s7: number;
	    momeryb7: number;
	    countb7: number;
	
	    static createFrom(source: any = {}) {
	        return new RedisMomery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.momeryForever = source["momeryForever"];
	        this.countForever = source["countForever"];
	        this.momerys3 = source["momerys3"];
	        this.counts3 = source["counts3"];
	        this.momeryb3s7 = source["momeryb3s7"];
	        this.countb3s7 = source["countb3s7"];
	        this.momeryb7 = source["momeryb7"];
	        this.countb7 = source["countb7"];
	    }
	}

}

export namespace utils {
	
	export class Connection {
	    key: string;
	    name: string;
	    host: string;
	    user: string;
	    password: string;
	    port: string;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.user = source["user"];
	        this.password = source["password"];
	        this.port = source["port"];
	    }
	}

}

