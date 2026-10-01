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
    if (name == "new") return 0;
    if (name == "add") return 2;
    if (name == "reserve") return 3;
    if (name == "cancel") return 1;
    if (name == "confirm") return 1;
    if (name == "available") return 1;
    throw invalid_argument("unknown call " + name);
}

// Reads T cases from stdin. Each case is k calls: "new", then "add item qty", "reserve order item qty", "cancel order", "confirm order" or "available item".
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
            Inventory* obj = nullptr;
            for (int i = 0; i < k; i++) {
                if (i > 0) line << ' ';
                if (names[i] == "new") {
                    delete obj;
                    obj = new Inventory();
                    line << "null";
                } else if (names[i] == "add") {
                    obj->add(params[i][0], params[i][1]);
                    line << "null";
                } else if (names[i] == "reserve") {
                    line << (obj->reserve(params[i][0], params[i][1], params[i][2]) ? "true" : "false");
                } else if (names[i] == "cancel") {
                    line << (obj->cancel(params[i][0]) ? "true" : "false");
                } else if (names[i] == "confirm") {
                    line << (obj->confirm(params[i][0]) ? "true" : "false");
                } else if (names[i] == "available") {
                    line << obj->available(params[i][0]);
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
