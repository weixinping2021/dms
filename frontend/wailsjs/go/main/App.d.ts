// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function AddCon(arg1:main.Connection):Promise<string>;

export function AnalyseRdb(arg1:string):Promise<string>;

export function DeleteCon(arg1:main.Connection):Promise<string>;

export function GetCons():Promise<Array<{[key: string]: string}>>;

export function GetConsStatus(arg1:string):Promise<{[key: string]: Array<main.MysqlConsStatus>}>;

export function GetConspercent(arg1:string):Promise<{[key: string]: string}>;

export function GetFullCons():Promise<Array<main.Connection>>;

export function GetMysqlLock(arg1:string):Promise<Array<main.MysqlProcessF>>;

export function GetMysqlProcesslist(arg1:string,arg2:string):Promise<Array<main.MysqlProcessF>>;

export function GetPeople():Promise<Array<main.Person>>;

export function GetPrefixkeys(arg1:string,arg2:string):Promise<Array<main.RedisKey>>;

export function GetRdbResultTitle():Promise<Array<{[key: string]: string}>>;

export function GetRedisKeys(arg1:string,arg2:string):Promise<Array<main.RedisKey>>;

export function GetRedisMemory(arg1:string):Promise<main.RedisMomery>;

export function GetRedisTop500Prefix(arg1:string):Promise<Array<main.RedisKey>>;

export function KillMysqlProcesses(arg1:string,arg2:Array<main.MysqlProcessF>):Promise<string>;

export function OpenDialog():Promise<string>;
