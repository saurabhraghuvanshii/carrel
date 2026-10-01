#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

// Reads T cases from stdin and prints one line per case. Each case is k calls: "new", "push x", "pop", "top" or "min".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
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
        vector<int> pushed(k);
        for (int i = 0; i < k; i++) {
            cin >> names[i];
            if (names[i] == "push") cin >> pushed[i];
        }
        try {
            ostringstream line;
            MinStack* stack = nullptr;
            for (int i = 0; i < k; i++) {
                if (i > 0) line << ' ';
                if (names[i] == "new") {
                    delete stack;
                    stack = new MinStack();
                    line << "null";
                } else if (names[i] == "push") {
                    stack->push(pushed[i]);
                    line << "null";
                } else if (names[i] == "pop") {
                    stack->pop();
                    line << "null";
                } else if (names[i] == "top") {
                    line << stack->top();
                } else {
                    line << stack->minimum();
                }
            }
            delete stack;
            out << line.str() << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
