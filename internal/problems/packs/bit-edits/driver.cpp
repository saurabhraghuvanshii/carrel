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

vector<int> readArray() {
    int n;
    cin >> n;
    vector<int> a(n);
    for (auto& x : a) cin >> x;
    return a;
}

string show(const vector<int>& a) {
    ostringstream o;
    o << '[';
    for (size_t i = 0; i < a.size(); i++) {
        if (i > 0) o << ", ";
        o << a[i];
    }
    o << ']';
    return o.str();
}

// Reads T cases from stdin and prints one line per case. Each case is "n q", then q operations: "get i", "set i", "clear i" or "update i b".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int n, q;
        cin >> n >> q;
        vector<string> ops(q);
        vector<int> at(q), bit(q, 0);
        for (int k = 0; k < q; k++) {
            cin >> ops[k] >> at[k];
            if (ops[k] == "update") cin >> bit[k];
        }
        try {
            vector<int> results(q);
            for (int k = 0; k < q; k++) {
                if (ops[k] == "get") {
                    results[k] = getBit(n, at[k]);
                    continue;
                }
                if (ops[k] == "set") n = setBit(n, at[k]);
                else if (ops[k] == "clear") n = clearBit(n, at[k]);
                else n = updateBit(n, at[k], bit[k]);
                results[k] = n;
            }
            out << show(results) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
