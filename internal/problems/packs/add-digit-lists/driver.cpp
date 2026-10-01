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

struct ListNode {
    int val;
    ListNode* next;
    ListNode(int v) : val(v), next(nullptr) {}
};

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

ListNode* build(const vector<int>& vals) {
    ListNode dummy(0);
    ListNode* tail = &dummy;
    for (int v : vals) {
        tail->next = new ListNode(v);
        tail = tail->next;
    }
    return dummy.next;
}

// Prints a list as [a, b, c]. More than limit nodes means a loop.
string showList(ListNode* head, size_t limit) {
    ostringstream o;
    o << '[';
    size_t count = 0;
    for (ListNode* p = head; p != nullptr; p = p->next) {
        if (++count > limit) throw runtime_error("the returned list has more nodes than it should. Is there a loop?");
        if (count > 1) o << ", ";
        o << p->val;
    }
    o << ']';
    return o.str();
}

// Reads T cases from stdin and prints one line per case. Each case is two lists of digits, lowest digit first.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        vector<int> a = readArray();
        vector<int> b = readArray();
        try {
            ListNode* r = addDigits(build(a), build(b));
            out << showList(r, max(a.size(), b.size()) + 1) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
