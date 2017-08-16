export class BCResolver {
    // All numbers to be stored as arrays, shown as numbers.
    // Represent as position:number

    private static get N() { return 4 };
    private allUniques;
    private workSet;

    constructor() {
        this.allUniques = this.genAll(BCResolver.N);
    }

    public genAll(n) {
        var set = [];
        var max = parseFloat(new Array(n + 1).join('9'));
        var num;
        var isUnique;
        for (var i = 0; i <= max; i++) {
            num = ('0000000000' + i).slice(-n);
            isUnique = true;
            for (var j = 0; j < n - 1; j++) {
                for (var k = j + 1; k < n; k++) {
                    if (num[j] == num[k]) isUnique = false;
                }
            }
            if (isUnique)
                set.push(num);
        }

        return set;
    }

    public searchInHash(hash, pos, val) {
        var hashLen = hash.length;
        for (var i = 0; i < hashLen; i++)
            if (hash[i].pos == pos && hash[i].val == val)
                return hash[i]

        return undefined;
    }

    public getHashTable(s) {
        var hash = [];
        var _this = this;
        s.forEach(function (num) {
            for (var i = 0; i < 4; i++) {
                // { pos: i, val: num[i], count: 0 }
                var res = _this.searchInHash(hash, i, num[i])
                if (res) res.count++
                else hash.push({ pos: i, val: num[i], count: 1 })
            }
        });

        hash.sort(function (a, b) {
            if (a.count > b.count) return 1
            if (a.count < b.count) return -1
            return 0;
        })

        return hash;
    }

    public printHash(hash) {
        var res = '';
        hash.forEach(function (h) {
            res += JSON.stringify(h) + '\n';
        })
        console.log(res);
    }

    public generateGuess(s, repeats?) {
        if (s.length == 1)
            return s;

        if (s.length < 100) {
            var uniq = this.findUniqueGuess(s)
            if (uniq && uniq.length) {
                return uniq[0];
            }
        }

        var sortedHash = this.getHashTable(s);
        var num = [];

        var i = 0;
        if (repeats) i++;

        while (num.length !== BCResolver.N) {
            if (num.indexOf(sortedHash[i].val) == -1)
                num.push(sortedHash[i].val)
            i++;
        }

        return num;
    }


    public respondToNum(num, guess) { // num is string, guess is number[]
        var response = { bulls: 0, cows: 0 };
        guess.forEach(function (dig, i) {
            dig = parseFloat(dig) + ""
            if (num.indexOf(dig) !== -1) {
                if (num.indexOf(dig) == i)
                    response.bulls++
                else
                    response.cows++
            }
        });

        return response;
    }

    public pruneSet(set, guess, bulls, cows) {
        let _this = this;
        var pruned = [];
        set.forEach(function (num, pos) {
            var numRes = _this.respondToNum(num, guess)
            if (numRes.bulls == bulls && numRes.cows == cows) {
                pruned.push(num)
            }
        });
        return pruned;
    }

    public findUniqueGuess(set) {
        if (set.length == 0) return set[0];

        var responses = [];
        var uniques = [];
        let _this = this;

        _this.allUniques.forEach(function (number) {
            var num = number.split('');

            set.forEach(function (n) {
                var res = _this.respondToNum(n, num);
                var resStr = 'b' + res.bulls + 'c' + res.cows;
                responses.push(resStr);
            });


            var matches = -responses.length;
            responses.forEach(function (res) {
                responses.forEach(function (r) {
                    if (res == r) matches++
                });
            })
            if (matches == 0) {
                uniques.push(num);
                return uniques;
            }

            responses = [];
        });
        return uniques;
    }

    public strNumToArray(number) {
        var num = [];
        number.split('').forEach(function (n) {
            num.push(parseInt(n, 10))
        });
        return num
    }

    public checkAll() {
        var num, times;
        var histogram = {};
        this.allUniques.forEach(function (number) {
            num = this.strNumToArray(number);
            times = this.playSingle(num)

            histogram[times] = histogram[times] ? histogram[times] + 1 : 1;
        })
        console.log(histogram);
    }

    public makeGuess() {
        this.workSet = this.allUniques.slice();
        var guess;
        var guesses = [];

        guess = this.generateGuess(this.workSet);
        if (guesses.indexOf(guess.join('')) !== -1) {
            guess = this.generateGuess(this.workSet, true);
        }

        guesses.push(guess.join(''));

        return guess.join('');
    }

    public prune(guess, bulls, cows) {
        this.allUniques = this.pruneSet(this.workSet, this.strNumToArray(guess), bulls, cows);
    }
}