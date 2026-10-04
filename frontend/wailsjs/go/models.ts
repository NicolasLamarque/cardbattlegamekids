export namespace store {
	
	export class Attack {
	    id: number;
	    characterId: number;
	    name: string;
	    color: string;
	    forceBonus: number;
	    manaCost: number;
	    sortOrder: number;
	
	    static createFrom(source: any = {}) {
	        return new Attack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.characterId = source["characterId"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.forceBonus = source["forceBonus"];
	        this.manaCost = source["manaCost"];
	        this.sortOrder = source["sortOrder"];
	    }
	}
	export class Character {
	    id: number;
	    nameFull: string;
	    nameNative: string;
	    imageLarge: string;
	    imageMedium: string;
	    imageLocal: string;
	    favourites: number;
	    role: string;
	    seriesId: number;
	    seriesTitle: string;
	    seriesColor: string;
	    force: number;
	    pv: number;
	    isCustomStats: boolean;
	    customData: string;
	    maxMana: number;
	    manaRegen: number;
	    maxAttacksOverride?: number;
	    maxAttacks: number;
	
	    static createFrom(source: any = {}) {
	        return new Character(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nameFull = source["nameFull"];
	        this.nameNative = source["nameNative"];
	        this.imageLarge = source["imageLarge"];
	        this.imageMedium = source["imageMedium"];
	        this.imageLocal = source["imageLocal"];
	        this.favourites = source["favourites"];
	        this.role = source["role"];
	        this.seriesId = source["seriesId"];
	        this.seriesTitle = source["seriesTitle"];
	        this.seriesColor = source["seriesColor"];
	        this.force = source["force"];
	        this.pv = source["pv"];
	        this.isCustomStats = source["isCustomStats"];
	        this.customData = source["customData"];
	        this.maxMana = source["maxMana"];
	        this.manaRegen = source["manaRegen"];
	        this.maxAttacksOverride = source["maxAttacksOverride"];
	        this.maxAttacks = source["maxAttacks"];
	    }
	}
	export class DraftModifier {
	    id: string;
	    label: string;
	    type: string;
	    kind: string;
	    amount: number;
	    rare: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DraftModifier(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.kind = source["kind"];
	        this.amount = source["amount"];
	        this.rare = source["rare"];
	    }
	}
	export class Player {
	    id: number;
	    name: string;
	    walletBalance: number;
	    gamesWon: number;
	    gamesPlayed: number;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.walletBalance = source["walletBalance"];
	        this.gamesWon = source["gamesWon"];
	        this.gamesPlayed = source["gamesPlayed"];
	    }
	}

}

