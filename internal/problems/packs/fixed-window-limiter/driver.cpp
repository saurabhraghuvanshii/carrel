#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iomanip>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stdexcept>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

int argCount(const string& name) {
    if (name == "new") return 2;
    if (name == "allow") return 2;
    throw invalid_argument("unknown call " + name);
}

// Reads T cases from stdin. Each case is k calls: "new limit window", then "allow user time".
// Prints one line per case with one result per call, null for calls that
// return nothing. Flushes after every case so a crash cannot hide finished
// results. Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int k;
        cin >> k;
        vector<string> names(k);
        vector<vector<int>> params(k);
        for (int i = 0; i < k; i++) {
            cin >> names[i];
            params[i].resize(argCount(names[i]));
            for (auto& p : params[i]) cin >> p;
        }
        ostringstream line;
        try {
            RateLimiter* obj = nullptr;
            for (int i = 0; i < k; i++) {
                if (i > 0) line << ' ';
                if (names[i] == "new") {
                    delete obj;
                    obj = new RateLimiter(params[i][0], params[i][1]);
                    line << "null";
                } else if (names[i] == "allow") {
                    line << (obj->allow(params[i][0], params[i][1]) ? "true" : "false");
                }
            }
            delete obj;
            out << line.str() << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
